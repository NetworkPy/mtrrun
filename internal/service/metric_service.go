package service

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/mtrrun/internal/model"
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
		log.Printf("metric with type=gauge and name=%s not found\n", name)

		return result, err
	}

	log.Printf("metric with type=gauge and name=%s found\n", name)

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
		log.Printf("metric with type=counter and name=%s not found\n", name)

		return result, err
	}

	log.Printf("metric with type=counter and name=%s found\n", name)

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
		log.Printf("metric with type=gauge and name=%s not found\n", dto.Name)
		log.Printf("metric with type=gauge and name=%s will be created\n", dto.Name)

		// Creating new metric if it not exists
		err = s.metRepo.InsertGauge(ctx, model.Gauge(dto))

		if err != nil {
			log.Printf("metric with type=gauge and name=%s was not created\n", dto.Name)

			return err
		}

		log.Printf("metric with type=gauge and name=%s created\n", dto.Name)

		return nil
	}

	log.Printf("metric with type=gauge and name=%s found\n", dto.Name)

	log.Printf("metric with type=gauge and name=%s will be updated\n", dto.Name)

	// Updating values if metric exists
	data.Value = dto.Value
	err = s.metRepo.UpdateGauge(ctx, data)

	if err != nil {
		log.Printf("metric with type=gauge and name=%s was not updated\n", dto.Name)

		return err
	}

	log.Printf("metric with type=gauge and name=%s updated\n", dto.Name)

	return nil
}

// PutCounter checking if metric exists.
// If yes then updating values, else creating new
func (s *MetricService) PutCounter(ctx context.Context, dto model.PutCounterDTO) error {
	// TODO: will make check sql.ErrNoRows when will implement repository with real database
	data, err := s.metRepo.SelectCounterByName(ctx, dto.Name)

	// Now error doesn't return fatal error. Any error = metric not exist
	if err != nil {
		log.Printf("metric with type=counter and name=%s not found\n", dto.Name)

		log.Printf("metric with type=counter and name=%s will be created\n", dto.Name)

		// Creating new metric if it not exists
		err = s.metRepo.InsertCounter(ctx, model.Counter(dto))

		if err != nil {
			log.Printf("metric with type=counter and name=%s was not created\n", dto.Name)

			return err
		}

		log.Printf("metric with type=counter and name=%s created\n", dto.Name)

		return nil
	}

	log.Printf("metric with type=counter and name=%s found\n", dto.Name)

	log.Printf("metric with type=counter and name=%s will be updated\n", dto.Name)

	// Adding new value to older
	data.Value += dto.Value
	err = s.metRepo.UpdateCounter(ctx, data)

	if err != nil {
		log.Printf("metric with type=counter and name=%s was not updated\n", dto.Name)

		return err
	}

	log.Printf("metric with type=counter and name=%s updated\n", dto.Name)

	return nil
}

// GetAll return all metrics. Calling repository methods for select all gauges and counters
func (s *MetricService) GetAll(ctx context.Context) ([]model.GetAllDTO, error) {
	dataGauge, err := s.metRepo.SelectGauge(ctx)

	if err != nil {
		log.Println("unable to find all gauge metrics")

		return nil, err
	}

	dataCounter, err := s.metRepo.SelectCounter(ctx)

	if err != nil {
		log.Println("unable to find all counter metrics")

		return nil, err
	}

	result := make([]model.GetAllDTO, 0)

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
