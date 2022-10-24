package model

// Gauge struct for data layer
type Gauge struct {
	Name  string
	Value float64
}

// Counter struct for data layer
type Counter struct {
	Name  string
	Value int64
}
