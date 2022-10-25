package service

import (
	"context"
	"fmt"
	"github.com/mtrrun/internal/model"
	"log"
	"strconv"
)

// metricRepository contract for repository layer
type metricRepository interface {
	SelectGaugeByName(ctx context.Context, name string) (model.Gauge, error)
	SelectCounterByName(ctx context.Context, name string) (model.Counter, error)
	SelectGauge(ctx context.Context) ([]model.Gauge, error)
	SelectCounter(ctx context.Context) ([]model.Counter, error)
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

// MetricServiceConfig config for MetricService
type MetricServiceConfig struct {
	MetRepo metricRepository
}

// NewMetricService constructor for MetricService
func NewMetricService(c *MetricServiceConfig) *MetricService {
	return &MetricService{
		metRepo: c.MetRepo,
	}
}

// GetGauge calling data layer and returning gauge metric or error
func (s *MetricService) GetGauge(ctx context.Context, name string) (model.GetGaugeDTO, error) {
	var result model.GetGaugeDTO

	metric, err := s.metRepo.SelectGaugeByName(ctx, name)

	if err != nil {
		log.Println("metric with type 'gauge' not found")

		return result, err
	}

	log.Println("metric with type 'gauge' found")

	// Mapping parameters to DTO
	result.Name = metric.Name
	result.Value = metric.Value

	return result, nil
}

// GetCounter calling data layer and returning counter metric or error
func (s *MetricService) GetCounter(ctx context.Context, name string) (model.GetCounterDTO, error) {
	var result model.GetCounterDTO

	metric, err := s.metRepo.SelectCounterByName(ctx, name)

	if err != nil {
		log.Println("metric with type 'counter' not found")

		return result, err
	}

	log.Println("metric with type 'counter' found")

	// Mapping parameters to DTO
	result.Name = metric.Name
	result.Value = metric.Value

	return result, nil
}

// PutGauge checking if metric exists.
// If yes then updating values, else creating new
func (s *MetricService) PutGauge(ctx context.Context, dto model.PutGaugeDTO) error {
	// TODO: will make check sql.ErrNoRows when will implement repository with real database
	data, err := s.metRepo.SelectGaugeByName(ctx, dto.Name)

	// Now error doesn't return fatal error. Any error = metric not exist
	if err != nil {
		log.Println("metric with type 'gauge' not found")
		log.Println("metric with type 'gauge' will be create")

		// Creating new metric if it not exists
		err = s.metRepo.InsertGauge(ctx, model.Gauge(dto))

		if err != nil {
			log.Println("metric with type 'gauge' didn't be create")

			return err
		}

		log.Println("metric with type 'gauge' created")

		return nil
	}

	log.Println("metric with type 'gauge' found")

	log.Println("metric with type 'gauge' will be update")

	// Updating values if metric exists
	data.Value = dto.Value
	err = s.metRepo.UpdateGauge(ctx, data)

	if err != nil {
		log.Println("metric with type 'gauge' didn't be update")

		return err
	}

	log.Println("metric with type 'gauge' updated")

	return nil
}

// PutCounter checking if metric exists.
// If yes then updating values, else creating new
func (s *MetricService) PutCounter(ctx context.Context, dto model.PutCounterDTO) error {
	// TODO: will make check sql.ErrNoRows when will implement repository with real database
	data, err := s.metRepo.SelectCounterByName(ctx, dto.Name)

	// Now error doesn't return fatal error. Any error = metric not exist
	if err != nil {
		log.Println("metric with type 'counter' not found")

		log.Println("metric with type 'counter' will be create")

		// Creating new metric if it not exists
		err = s.metRepo.InsertCounter(ctx, model.Counter(dto))

		if err != nil {
			log.Println("metric with type 'counter' didn't be create")

			return err
		}

		log.Println("metric with type 'counter' created")

		return nil
	}

	log.Println("metric with type 'counter' found")

	log.Println("metric with type 'counter' will be update")

	// Adding new value to older
	data.Value += dto.Value
	err = s.metRepo.UpdateCounter(ctx, data)

	if err != nil {
		log.Println("metric with type 'counter' didn't be update")

		return err
	}

	log.Println("metric with type 'counter' updated")

	return nil
}

// GetAll return all metrics. Calling repository methods for select all gauges and counters
func (s *MetricService) GetAll(ctx context.Context) ([]model.GetAllDTO, error) {
	result := make([]model.GetAllDTO, 0)

	dataGauge, err := s.metRepo.SelectGauge(ctx)

	if err != nil {
		log.Println("unable to find all gauge metrics")

		return result, err
	}

	dataCounter, err := s.metRepo.SelectCounter(ctx)

	if err != nil {
		log.Println("unable to find all counter metrics")

		return result, err
	}

	for i := 0; i < len(dataGauge); i++ {
		result = append(result, model.GetAllDTO{
			Name:  dataGauge[i].Name,
			Value: strconv.FormatFloat(dataGauge[i].Value, 'f', -1, 64),
		})
	}

	for i := 0; i < len(dataCounter); i++ {
		result = append(result, model.GetAllDTO{
			Name:  dataCounter[i].Name,
			Value: fmt.Sprintf("%d", dataCounter[i].Value),
		})
	}

	return result, nil
}
