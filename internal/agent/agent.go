package agent

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Metric is interface for
// all metric types.
type Metric interface {
	Desc() Description
	GetValue() string
}

// Tracker is interface which declarate methods for
// all container realisations with metrics.
// Abstraction layer which controls metrics.
// Doing create, delete, values update and other operations
type Tracker interface {
	Track(metric Metric)
	Untrack(metric Metric)
	Status() []Status
}

// Counter is very simple realisation from Prometheus library
type Counter interface {
	Metric
	Inc()
}

// Gauge is analog from Prometheus library
// for calculate variable values
type Gauge interface {
	Metric
	Set(float64)
	Inc()
	Dec()
	Add(float64)
	Sub(float64)
}

type Histogram interface {
	Metric
	// TODO: do it in next time
}

type Summary interface {
	Metric
	// TODO: do it in next time
}

type Logger interface {
	Info()
	Debug()
	Error()
}

// Client interface for client that sending metrics to other service
type Client interface {
	DoRequest(method string, url string, header map[string]string, body []byte) error
	Shutdown()
}

// Description for metrics
type Description struct {
	Name string
	Help string
}

// Status information about metrics
type Status struct {
	Name       string
	MetricType string
	Value      string
}

const (
	defaultReportInterval       = 2
	defaultPollInterval         = 10
	defaultMaxRequestsPerMoment = 5

	contentTypeHeader  = "Content-Type"
	defaultContentType = "text/plain"

	gaugeType   = "gauge"
	counterType = "counter"
	unknownType = "unknown"
)

// Agent calling container for get, create, update, delete metrics.
// Implementation Facade pattern
type Agent struct {

	// Container with metrics
	container Tracker

	// Client for sending request with metrics to other services
	client Client

	// Time interval in seconds for sending request to other service.
	// If reportInterval is empty that will be use default value - 2 second.
	reportInterval time.Duration

	// Interval in seconds which call update metrics.
	// If pollInterval is empty that will be use default value - 10 second.
	pollInterval time.Duration

	// Channel and sync.Once for gracefully shutdown
	exit       chan struct{}
	onceCloser sync.Once

	host                 string
	maxRequestsPerMoment int
}

// Config configuration list for Agent
type Config struct {
	//Container      Tracker
	//Client         Client
	ReportInterval time.Duration
	PollInterval   time.Duration

	Host                 string
	MaxRequestsPerMoment int

	Timeout      time.Duration // Time in seconds
	MaxIdleConns int           // Max cached connections
}

// New constructor for Agent
func New(c *Config) (*Agent, error) {
	if c.ReportInterval <= 0 {
		c.ReportInterval = defaultReportInterval
	}

	if c.PollInterval <= 0 {
		c.PollInterval = defaultPollInterval
	}

	if c.MaxRequestsPerMoment <= 0 {
		c.MaxRequestsPerMoment = defaultMaxRequestsPerMoment
	}

	return &Agent{
		container:      NewTracker(),
		client:         NewClient(c.Timeout, c.MaxIdleConns),
		reportInterval: c.ReportInterval,
		pollInterval:   c.PollInterval,

		host:                 c.Host,
		maxRequestsPerMoment: c.MaxRequestsPerMoment,
	}, nil
}

// Track adding metric to track
func (a *Agent) Track(metric Metric) {
	a.container.Track(metric)
}

// Untrack removing metric
func (a *Agent) Untrack(metric Metric) {
	a.container.Untrack(metric)
}

// Status returning actual metrics state
func (a *Agent) Status() []Status {
	return a.container.Status()
}

func (a *Agent) CustomTracker(t Tracker) {
	a.container = t
}

// Run call unblocking operation and start event
// cycle with collect metrics
// from runtime
func (a *Agent) Run() {
	reportTicker := time.NewTicker(a.reportInterval)
	a.exit = make(chan struct{})

	for {
		select {
		case <-a.exit:
			// Gracefully shutdown ticker
			reportTicker.Stop()
			log.Printf("agent been gracefully shutdown")

			return
		case <-reportTicker.C:
			a.report()
		}
	}
}

// Shutdown call functional for correct exits from program
func (a *Agent) Shutdown() {
	// Gracefully shutdown all agent's components
	a.onceCloser.Do(func() {
		close(a.exit)
		a.client.Shutdown()
	})
}

// Sending report with metrics
func (a *Agent) report() {
	s := a.container.Status()

	// For implementation semaphore
	channel := make(chan struct{}, a.maxRequestsPerMoment)

	var wg sync.WaitGroup

	for i := 0; i < len(s); i++ {

		wg.Add(1)
		go func(x int) {
			defer wg.Done()

			channel <- struct{}{}

			url := fmt.Sprintf("http://%s/update/%s/%s/%s", a.host, s[x].MetricType, s[x].Name, s[x].Value)

			log.Printf("start of request to url: %s\n", url)

			err := a.client.DoRequest(http.MethodPost, url, map[string]string{contentTypeHeader: defaultContentType}, nil)

			if err != nil {
				log.Printf("request ended with error: %s\n", err)
			} else {
				log.Printf("request ended without error")
			}

			<-channel
		}(i)
	}

	wg.Wait()
}

func getMetricType(met Metric) string {
	switch met.(type) {
	case Gauge:
		return gaugeType
	case Counter:
		return counterType
	default:
		return unknownType
	}
}
