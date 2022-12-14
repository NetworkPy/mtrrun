package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/mtrrun/internal/agent"
	"github.com/mtrrun/internal/config"
)

// Golang runtime metrics:
// Metric name: "Alloc", type: gauge
// Metric name: "BuckHashSys", type: gauge
// Metric name: "Frees", type: gauge
// Metric name: "GCCPUFraction", type: gauge
// Metric name: "GCSys", type: gauge
// Metric name: "HeapAlloc", type: gauge
// Metric name: "HeapIdle", type: gauge
// Metric name: "HeapInuse", type: gauge
// Metric name: "HeapObjects", type: gauge
// Metric name: "HeapReleased", type: gauge
// Metric name: "HeapSys", type: gauge
// Metric name: "LastGC", type: gauge
// Metric name: "Lookups", type: gauge
// Metric name: "MCacheInuse", type: gauge
// Metric name: "MCacheSys", type: gauge
// Metric name: "MSpanInuse", type: gauge
// Metric name: "MSpanSys", type: gauge
// Metric name: "Mallocs", type: gauge
// Metric name: "NextGC", type: gauge
// Metric name: "NumForcedGC", type: gauge
// Metric name: "NumGC", type: gauge
// Metric name: "OtherSys", type: gauge
// Metric name: "PauseTotalNs", type: gauge
// Metric name: "StackInuse", type: gauge
// Metric name: "StackSys", type: gauge
// Metric name: "Sys", type: gauge
// Metric name: "TotalAlloc", type: gauge

// Other custom metrics:
// Metric name: "PollCount", type: counter — counter, incremented by 1 each time a metric from the runtime package is updated.
// Metric name: "RandomValue", type: gauge — random value.

const (
	MetricAlloc       = "Alloc"
	MetricBuckHashSys = "BuckHashSys"
	MetricFrees       = "Frees"
	MetricGCCPUFracti = "GCCPUFracti"
	MetricGCSys       = "GCSys"
	MetricHeapAlloc   = "HeapAlloc"
	MetricHeapIdle    = "HeapIdle"
	MetricHeapInuse   = "HeapInuse"
	MetricHeapObjects = "HeapObjects"
	MetricHeapRelease = "HeapRelease"
	MetricHeapSys     = "HeapSys"
	MetricLastGC      = "LastGC"
	MetricLookups     = "Lookups"
	MetricMCacheInuse = "MCacheSys"
	MetricMCacheSys   = "CacheSys"
	MetricMSpanInuse  = "MSpanInuse"
	MetricMSpanSys    = "MSpanSys"
	MetricMallocs     = "Mallocs"
	MetricNextGC      = "NextGC"
	MetricNumForcedGC = "NumForcedGC"
	MetricNumGC       = "NumGC"
	MetricOtherSys    = "OtherSys"
	MetricPauseTotalN = "PauseTotalN"
	MetricStackInuse  = "StackInuse"
	MetricStackSys    = "StackSys"
	MetricSys         = "Sys"
	MetricTotalAlloc  = "TotalAlloc"
	MetricPollCount   = "PollCount"
	MetricRandomValue = "RandomValue"
)

// For configuration
const (
	defaultHost                 = "127.0.0.1:8080"
	defaultTimeout              = 5 * time.Second
	defaultMaxIdleConns         = 5
	defaultMaxRequestsPerMoment = 5
	defaultReportInterval       = 10 * time.Second
	defaultPollInterval         = 2 * time.Second
)

func main() {
	//var path string
	//flag.StringVar(&path, "config", "config.yaml", "your name")
	//flag.Parse()
	//
	//c, err := config.ReadAgentConfig(path)

	//if err != nil {
	//	log.Fatalf(err.Error())
	//}

	// TODO: hardcode config because CI/CD don't load config file
	c := &config.AgentConfig{
		Host:                 defaultHost,
		Timeout:              defaultTimeout,
		MaxIdleConns:         defaultMaxIdleConns,
		MaxRequestsPerMoment: defaultMaxRequestsPerMoment,
		ReportInterval:       defaultReportInterval,
		PollInterval:         defaultPollInterval,
	}

	a, err := agent.New(&agent.Config{
		ReportInterval:       c.ReportInterval,
		PollInterval:         c.PollInterval,
		Host:                 c.Host,
		MaxRequestsPerMoment: c.MaxRequestsPerMoment,
		Timeout:              c.Timeout,
		MaxIdleConns:         c.MaxIdleConns,
	})
	if err != nil {
		log.Fatalf("failed to create agent: %s", err)
	}

	// Channel for gracefully shutdown
	exitCh := make(chan struct{})

	// Starting new goroutine inside initMetrics
	initMetrics(a, c.PollInterval, exitCh)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// Starting agent cycle
		a.Run()
	}()

	log.Println("agent started")

	<-done

	go func() {
		a.Shutdown()
		log.Println("agent stopped")
		close(exitCh)
	}()

	<-exitCh

	log.Println("agent exited properly")
}

// initMetrics init all metrics and start job with runtime.MemStats in a separate goroutine
func initMetrics(a *agent.Agent, pollInterval time.Duration, exitCh chan struct{}) {
	var m = &runtime.MemStats{}

	// Init runtime metrics
	metricAlloc := agent.NewGauge(MetricAlloc, "")
	metricBuckHashSys := agent.NewGauge(MetricBuckHashSys, "")
	metricFrees := agent.NewGauge(MetricFrees, "")
	metricGCCPUFracti := agent.NewGauge(MetricGCCPUFracti, "")
	metricGCSys := agent.NewGauge(MetricGCSys, "")
	metricHeapAlloc := agent.NewGauge(MetricHeapAlloc, "")
	metricHeapIdle := agent.NewGauge(MetricHeapIdle, "")
	metricHeapInuse := agent.NewGauge(MetricHeapInuse, "")
	metricHeapObjects := agent.NewGauge(MetricHeapObjects, "")
	metricHeapRelease := agent.NewGauge(MetricHeapRelease, "")
	metricHeapSys := agent.NewGauge(MetricHeapSys, "")
	metricLastGC := agent.NewGauge(MetricLastGC, "")
	metricLookups := agent.NewGauge(MetricLookups, "")
	metricMCacheInuse := agent.NewGauge(MetricMCacheInuse, "")
	metricMCacheSys := agent.NewGauge(MetricMCacheSys, "")
	metricMSpanInuse := agent.NewGauge(MetricMSpanInuse, "")
	metricMSpanSys := agent.NewGauge(MetricMSpanSys, "")
	metricMallocs := agent.NewGauge(MetricMallocs, "")
	metricNextGC := agent.NewGauge(MetricNextGC, "")
	metricNumForcedGC := agent.NewGauge(MetricNumForcedGC, "")
	metricNumGC := agent.NewGauge(MetricNumGC, "")
	metricOtherSys := agent.NewGauge(MetricOtherSys, "")
	metricPauseTotalN := agent.NewGauge(MetricPauseTotalN, "")
	metricStackInuse := agent.NewGauge(MetricStackInuse, "")
	metricStackSys := agent.NewGauge(MetricStackSys, "")
	metricSys := agent.NewGauge(MetricSys, "")
	metricTotalAlloc := agent.NewGauge(MetricTotalAlloc, "")

	// Init custom metrics
	metricPollCount := agent.NewCounter(MetricPollCount, "")
	metricRandomValue := agent.NewGauge(MetricRandomValue, "")

	// Adding all metrics to track
	a.Track(metricAlloc)
	a.Track(metricBuckHashSys)
	a.Track(metricFrees)
	a.Track(metricGCCPUFracti)
	a.Track(metricGCSys)
	a.Track(metricHeapAlloc)
	a.Track(metricHeapIdle)
	a.Track(metricHeapInuse)
	a.Track(metricHeapObjects)
	a.Track(metricHeapRelease)
	a.Track(metricHeapSys)
	a.Track(metricLastGC)
	a.Track(metricLookups)
	a.Track(metricMCacheInuse)
	a.Track(metricMCacheSys)
	a.Track(metricMSpanInuse)
	a.Track(metricMSpanSys)
	a.Track(metricMallocs)
	a.Track(metricNextGC)
	a.Track(metricNumForcedGC)
	a.Track(metricNumGC)
	a.Track(metricOtherSys)
	a.Track(metricPauseTotalN)
	a.Track(metricStackInuse)
	a.Track(metricStackSys)
	a.Track(metricSys)
	a.Track(metricTotalAlloc)
	a.Track(metricPollCount)
	a.Track(metricRandomValue)

	// Ticker for cycle with collecting metrics
	pollTicker := time.NewTicker(pollInterval)

	go func() {
		for {
			select {
			case <-exitCh:
				pollTicker.Stop()
				log.Println("main cycle with collecting metrics stopped")

				return
			case <-pollTicker.C:
				runtime.ReadMemStats(m)
				rand.Seed(time.Now().UTC().UnixNano())

				metricAlloc.Set(float64(m.Alloc))
				metricBuckHashSys.Set(float64(m.BuckHashSys))
				metricFrees.Set(float64(m.Frees))
				metricGCCPUFracti.Set(m.GCCPUFraction)
				metricGCSys.Set(float64(m.GCSys))
				metricHeapAlloc.Set(float64(m.HeapAlloc))
				metricHeapIdle.Set(float64(m.HeapIdle))
				metricHeapInuse.Set(float64(m.HeapInuse))
				metricHeapObjects.Set(float64(m.HeapObjects))
				metricHeapRelease.Set(float64(m.HeapReleased))
				metricHeapSys.Set(float64(m.HeapSys))
				metricLastGC.Set(float64(m.LastGC))
				metricLookups.Set(float64(m.Lookups))
				metricMCacheInuse.Set(float64(m.MCacheInuse))
				metricMCacheSys.Set(float64(m.MCacheSys))
				metricMSpanInuse.Set(float64(m.MSpanInuse))
				metricMSpanSys.Set(float64(m.MSpanSys))
				metricMallocs.Set(float64(m.Mallocs))
				metricNextGC.Set(float64(m.NextGC))
				metricNumForcedGC.Set(float64(m.NumForcedGC))
				metricNumGC.Set(float64(m.NumGC))
				metricOtherSys.Set(float64(m.OtherSys))
				metricPauseTotalN.Set(float64(m.PauseTotalNs))
				metricStackInuse.Set(float64(m.StackInuse))
				metricStackSys.Set(float64(m.StackSys))
				metricSys.Set(float64(m.Sys))
				metricTotalAlloc.Set(float64(m.TotalAlloc))

				metricPollCount.Inc()
				metricRandomValue.Set(float64(rand.Int63()))
			}
		}
	}()
}
