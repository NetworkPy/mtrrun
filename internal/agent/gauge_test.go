package agent

import (
	"sync"
	"testing"
)

func TestGauge_With(t *testing.T) {
	runtimeGauge := &gauge{
		mu:  sync.RWMutex{},
		val: 0,
		d: &Description{
			Name: "test",
			Help: "test",
		},
	}

	wg := sync.WaitGroup{}

	for i := 0; i < 900; i++ {
		wg.Add(1)
		go func() {
			runtimeGauge.Inc()
			defer wg.Done()
		}()
	}

	wg.Wait()
}
