package service

import (
	"context"
	"fmt"
	"github.com/mtrrun/internal/model"
)

const (
	gauge   = "gauge"
	counter = "counter"
)

// metricRepository contract for repository layer
type metricRepository interface {
	SelectGaugeByName(name string) (model.Gauge, error)
	SelectCounterByName(name string) (model.Counter, error)
	InsertGauge(ctx context.Context, metric model.Gauge) error
	InsertCounter(ctx context.Context, metric model.Counter) error
	UpdateGauge(ctx context.Context, curr model.Gauge) error
	UpdateCounter(ctx context.Context, curr model.Counter) error
	DeleteGauge(ctx context.Context, name string) error
	DeleteCounter(ctx context.Context, name string) error
}

// MetricService layer with business logic for metrics
type MetricService struct {
	metRepo metricRepository
}

// CheckMetricType for check metric type
func (s *MetricService) CheckMetricType(t string) error {
	switch t {
	case gauge:

		return nil
	case counter:

		return nil
	}

	return fmt.Errorf("unknown metric type. Expected %s or %s. Actual: %s", gauge, counter, t)
}
