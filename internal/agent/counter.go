package agent

import (
	"fmt"
	"sync"
)

// Implementing Counter interface
type counter struct {
	mu  sync.RWMutex
	val int64
	d   *Description
}

func (c *counter) Desc() Description {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return *c.d
}

// Inc increments the Counter by 1. Use Add to increment it by arbitrary values
func (c *counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.val++
}

// GetValue returned value
func (c *counter) GetValue() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return fmt.Sprintf("%d", c.val)
}

func NewCounter(name string, help string) Counter {
	return &counter{
		val: 0,
		d: &Description{
			Name: name,
			Help: help,
		},
	}
}
