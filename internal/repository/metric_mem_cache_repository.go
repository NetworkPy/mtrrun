package repository

import (
	"context"
	"fmt"
	"github.com/mtrrun/internal/model"
	"sync"
)

// MetricMemCache in memory cache for server with metrics.
// Contains gauge and counter types for metrics.
type MetricMemCache struct {
	gaugeMu   sync.RWMutex
	counterMu sync.RWMutex

	gauge   map[string]model.Gauge
	counter map[string]model.Counter
}

// NewMetricMemCache Constructor for MetricMemCache
func NewMetricMemCache() *MetricMemCache {
	return &MetricMemCache{
		gauge:   make(map[string]model.Gauge),
		counter: make(map[string]model.Counter),
	}
}

// SelectGaugeByName selecting gauge metric by name
func (c *MetricMemCache) SelectGaugeByName(name string) (model.Gauge, error) {
	c.gaugeMu.RLock()
	defer c.gaugeMu.RUnlock()

	var metric model.Gauge

	if _, ok := c.gauge[name]; ok {
		metric = c.gauge[name]

		return metric, nil
	}

	return metric, fmt.Errorf("gauge metric by name=%s not found", name)
}

// SelectCounterByName selecting counter metric by name
func (c *MetricMemCache) SelectCounterByName(name string) (model.Counter, error) {
	c.counterMu.RLock()
	defer c.counterMu.RUnlock()

	var metric model.Counter

	if _, ok := c.counter[name]; ok {
		metric = c.counter[name]

		return metric, nil
	}

	return metric, fmt.Errorf("counter metric by name=%s not found", name)
}

func (c *MetricMemCache) InsertGauge(ctx context.Context, metric model.Gauge) error {
	c.gaugeMu.Lock()
	defer c.gaugeMu.Unlock()

	c.gauge[metric.Name] = metric

	return nil
}

// InsertCounter inserting counter metric if it is not exist
func (c *MetricMemCache) InsertCounter(ctx context.Context, metric model.Counter) error {
	c.counterMu.Lock()
	defer c.counterMu.Unlock()

	c.counter[metric.Name] = metric

	return nil
}

// UpdateGauge updating counter metric. It is assumed that the metric exists
func (c *MetricMemCache) UpdateGauge(ctx context.Context, curr model.Gauge) error {
	c.gaugeMu.RLock()
	defer c.gaugeMu.RLock()

	// We can do c.gauge[curr.Name] = curr.
	// But in continue we could want to update
	// not all fields

	prev := c.gauge[curr.Name]

	prev.Value = curr.Value

	c.gauge[prev.Name] = prev

	return nil
}

// UpdateCounter updating counter metric. It is assumed that the metric exists
func (c *MetricMemCache) UpdateCounter(ctx context.Context, curr model.Counter) error {
	c.counterMu.RLock()
	defer c.counterMu.RLock()

	prev := c.counter[curr.Name]

	prev.Value = curr.Value

	c.counter[prev.Name] = prev

	return nil
}

// DeleteGauge deleting metric with gauge type
func (c *MetricMemCache) DeleteGauge(ctx context.Context, name string) error {
	c.gaugeMu.Lock()
	defer c.gaugeMu.Lock()

	delete(c.gauge, name)

	return nil
}

// DeleteCounter deleting metric with counter type
func (c *MetricMemCache) DeleteCounter(ctx context.Context, name string) error {
	c.counterMu.Lock()
	defer c.counterMu.Lock()

	delete(c.counter, name)

	return nil
}
