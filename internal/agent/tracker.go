package agent

import (
	"sync"
)

// tracker is main storage with metrics
// which added or removed metrics
type tracker struct {
	mu sync.RWMutex

	// storages with all metrics

	metrics map[string]Metric

	// TODO: histogram, summary
}

// Track added metric to list with all metrics
func (r *tracker) Track(met Metric) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.metrics[met.Desc().Name] = met
}

// Untrack removed metric from list with all metrics
func (r *tracker) Untrack(met Metric) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.metrics, met.Desc().Name)
}

// Status returned information about actual metrics state
func (r *tracker) Status() []Status {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s := make([]Status, len(r.metrics))

	count := 0
	for k, v := range r.metrics {
		s[count] = Status{
			Name:       k,
			MetricType: getMetricType(v),
			Value:      v.GetValue(),
		}

		count++
	}

	return s
}

// NewTracker constructor for Tracker
func NewTracker() Tracker {

	return &tracker{
		metrics: make(map[string]Metric),
	}
}
