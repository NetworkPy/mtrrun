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

// GetGaugeDTO data transfer object between
// handler layer and service layer for getting gauge
type GetGaugeDTO struct {
	Name  string
	Value float64
}

// GetCounterDTO data transfer object between
// handler layer and service layer for getting counter
type GetCounterDTO struct {
	Name  string
	Value int64
}

// PutGaugeDTO data transfer object between
// handler layer and service layer for putting gauge
type PutGaugeDTO struct {
	Name  string
	Value float64
}

// PutCounterDTO data transfer object between
// handler layer and service layer for putting counter
type PutCounterDTO struct {
	Name  string
	Value int64
}
