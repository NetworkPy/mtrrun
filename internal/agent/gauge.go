package agent

import (
	"fmt"
	"sync"
)

// Implementing Gauge interface
type gauge struct {
	mu  sync.RWMutex
	val float64
	d   *Description
}

func (g *gauge) Desc() Description {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return *g.d
}

// Set sets the Gauge to an arbitrary value.
func (g *gauge) Set(val float64) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.val = val
}

// Inc increments the Gauge by 1. Use Add to increment it by arbitrary values
func (g *gauge) Inc() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.val++
}

// Dec decrements the Gauge by 1. Use Sub to decrement it by arbitrary values
func (g *gauge) Dec() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.val--
}

// Add adds the given value to the Gauge
func (g *gauge) Add(val float64) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.val += val
}

// Sub subtracts the given value from the Gauge
func (g *gauge) Sub(val float64) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.val -= val
}

// GetValue returned value
func (g *gauge) GetValue() string {
	var val float64

	g.mu.Lock()

	val = g.val

	g.mu.Unlock()

	return fmt.Sprintf("%.2f", val)
}

func NewGauge(name string, help string) Gauge {
	return &gauge{
		val: 0,
		d: &Description{
			Name: name,
			Help: help,
		},
	}
}
