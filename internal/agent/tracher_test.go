package agent

import (
	"fmt"
	"testing"
)

func TestTracker_Status(t *testing.T) {

	c1 := NewCounter("best_counter_total", "best counter ever")
	c2 := NewGauge("best_gauge_total", "best gauge total")

	track := NewTracker()

	track.Track(c1)

	track.Track(c2)

	c1.Inc()
	c1.Inc()
	c2.Dec()

	fmt.Println(c2.GetValue())

	fmt.Println(track.Status())

}
