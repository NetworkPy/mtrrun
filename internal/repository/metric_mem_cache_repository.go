package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/mtrrun/internal/model"
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
func (c *MetricMemCache) SelectGaugeByName(ctx context.Context, name string) (model.Gauge, error) {
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
func (c *MetricMemCache) SelectCounterByName(ctx context.Context, name string) (model.Counter, error) {
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

	if _, ok := c.gauge[metric.Name]; ok {
		return fmt.Errorf("unable to create metric with name=%s and type=gauge. Metric exists", metric.Name)
	}

	c.gauge[metric.Name] = metric

	return nil
}

// InsertCounter inserting counter metric if it is not exist
func (c *MetricMemCache) InsertCounter(ctx context.Context, metric model.Counter) error {
	c.counterMu.Lock()
	defer c.counterMu.Unlock()

	if _, ok := c.counter[metric.Name]; ok {
		return fmt.Errorf("unable to create metric with name=%s and type=counter. Metric is exists", metric.Name)
	}

	c.counter[metric.Name] = metric

	return nil
}

// UpdateGauge updating counter metric. It is assumed that the metric exists
func (c *MetricMemCache) UpdateGauge(ctx context.Context, curr model.Gauge) error {
	c.gaugeMu.Lock()
	defer c.gaugeMu.Unlock()

	// We can do c.gauge[curr.Name] = curr.
	// But in continue we could want to update
	// not all fields

	if _, ok := c.gauge[curr.Name]; !ok {
		return fmt.Errorf("unable to update metric with name=%s and type=gauge. Metric is not exists", curr.Name)
	}

	prev := c.gauge[curr.Name]

	prev.Value = curr.Value

	c.gauge[prev.Name] = prev

	return nil
}

// UpdateCounter updating counter metric. It is assumed that the metric exists
func (c *MetricMemCache) UpdateCounter(ctx context.Context, curr model.Counter) error {
	c.counterMu.Lock()
	defer c.counterMu.Unlock()

	if _, ok := c.counter[curr.Name]; !ok {
		return fmt.Errorf("unable to update metric with name=%s and type=counter. Metric is not exists", curr.Name)
	}

	prev := c.counter[curr.Name]

	prev.Value = curr.Value

	c.counter[prev.Name] = prev

	return nil
}

// DeleteGauge deleting metric with gauge type
func (c *MetricMemCache) DeleteGauge(ctx context.Context, name string) error {
	c.gaugeMu.Lock()
	defer c.gaugeMu.Unlock()

	if _, ok := c.gauge[name]; !ok {
		return fmt.Errorf("unable to delete metric with name=%s and type=gauge. Metric is not exists", name)
	}

	delete(c.gauge, name)

	return nil
}

// DeleteCounter deleting metric with counter type
func (c *MetricMemCache) DeleteCounter(ctx context.Context, name string) error {
	c.counterMu.Lock()
	defer c.counterMu.Unlock()

	if _, ok := c.counter[name]; !ok {
		return fmt.Errorf("unable to delete metric with name=%s and type=counter. Metric is not exists", name)
	}

	delete(c.counter, name)

	return nil
}

// SelectGauge selecting all metrics with type gauge
func (c *MetricMemCache) SelectGauge(ctx context.Context) ([]model.Gauge, error) {
	c.gaugeMu.RLock()
	defer c.gaugeMu.RUnlock()

	result := make([]model.Gauge, 0, len(c.gauge))

	for _, v := range c.gauge {
		result = append(result, v)
	}

	return result, nil
}

// SelectCounter selecting all metrics with type counter
func (c *MetricMemCache) SelectCounter(ctx context.Context) ([]model.Counter, error) {
	c.counterMu.RLock()
	defer c.counterMu.RUnlock()

	result := make([]model.Counter, 0, len(c.counter))

	for _, v := range c.counter {
		result = append(result, v)
	}

	return result, nil
}
