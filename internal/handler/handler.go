package handler

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mtrrun/internal/model"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

const (
	gauge   = "gauge"
	counter = "counter"
)

type metricService interface {
	GetGauge(ctx context.Context, name string) (model.GetGaugeDTO, error)
	GetCounter(ctx context.Context, name string) (model.GetCounterDTO, error)
	GetAll(ctx context.Context) ([]model.GetAllDTO, error)
	PutGauge(ctx context.Context, dto model.PutGaugeDTO) error
	PutCounter(ctx context.Context, dto model.PutCounterDTO) error
}

// Handler implementing all handlers for server
type Handler struct {
	metSrv metricService
}

// Config for Handler
type Config struct {
	Router *mux.Router
	MetSrv metricService
}

// New is constructor for Handler
func New(c *Config) {
	h := Handler{
		metSrv: c.MetSrv,
	}

	c.Router.HandleFunc("/", h.GetStaticAllMetrics).Methods(http.MethodGet)
	c.Router.HandleFunc("/update/{metric_type}/{metric_name}/{value}", h.UpdateMetric).Methods(http.MethodPost)
	c.Router.HandleFunc("/value/{metric_type}/{metric_name}", h.GetMetric).Methods(http.MethodGet)
}

// UpdateMetric accepts request for create or update metrics
func (h *Handler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	metricType, ok := vars["metric_type"]

	if !ok {
		http.Error(w, "unable parse path parameter 'metric_type'", http.StatusBadRequest)

		return
	}

	metricName, ok := vars["metric_name"]

	if !ok {
		http.Error(w, "unable parse path parameter 'metric_name'", http.StatusBadRequest)

		return
	}

	value, ok := vars["value"]

	if !ok {
		http.Error(w, "unable parse path parameter 'value'", http.StatusBadRequest)

		return
	}

	// Validate path params

	if len(metricName) == 0 {
		log.Println("unable to parse name. Expected: string with length > 0")
		http.Error(w, "unable to parse name. Expected: string with length > 0", http.StatusBadRequest)

		return
	}

	switch metricType {
	case gauge:
		valueFloat64, err := strconv.ParseFloat(value, 64)

		if err != nil {
			log.Printf("unable to parse value. Expected: float. Actual: %s", value)
			http.Error(w, fmt.Sprintf("unable to parse value. Expected: float. Actual: %s", value), http.StatusBadRequest)

			return
		}

		err = h.metSrv.PutGauge(ctx, model.PutGaugeDTO{
			Name:  metricName,
			Value: valueFloat64,
		})

		if err != nil {
			log.Printf("unable to update/create gauge metric with name=%s and value=%s", metricName, value)
			http.Error(w, fmt.Sprintf("unable to update/create gauge metric with name=%s and value=%s",
				metricName, value), http.StatusInternalServerError)

			return
		}
	case counter:
		valueInt64, err := strconv.ParseInt(value, 10, 64)

		if err != nil {
			log.Printf("unable to parse value. Expected: int. Actual: %s", value)
			http.Error(w, fmt.Sprintf("unable to parse value. Expected: int. Actual: %s",
				value), http.StatusBadRequest)

			return
		}

		err = h.metSrv.PutCounter(ctx, model.PutCounterDTO{
			Name:  metricName,
			Value: valueInt64,
		})

		if err != nil {
			log.Printf("unable to update/create counter metric with name=%s and value=%s", metricName, value)
			http.Error(w, fmt.Sprintf("unable to update/create counter metric with name=%s and value=%s",
				metricName, value), http.StatusInternalServerError)

			return
		}
	default:
		log.Printf("unknown metric type. Expected %s or %s. Actual: %s\n", gauge, counter, metricType)
		http.Error(w, fmt.Sprintf("unknown metric type. Expected %s or %s. Actual: %s\n", gauge, counter, metricType), http.StatusNotImplemented)

		return
	}

	_, err := w.Write([]byte("OK"))

	if err != nil {
		w.WriteHeader(500)

		return
	}
}

// GetMetric return information about metric by name if it exists
func (h *Handler) GetMetric(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	metricType, ok := vars["metric_type"]

	if !ok {
		http.Error(w, "unable parse path parameter 'metric_type'", http.StatusBadRequest)

		return
	}

	metricName, ok := vars["metric_name"]

	if !ok {
		http.Error(w, "unable parse path parameter 'metric_name'", http.StatusBadRequest)

		return
	}

	if len(metricName) == 0 {
		http.Error(w, "unable parse path parameter 'metric_type'", http.StatusBadRequest)

		return
	}

	switch metricType {
	case gauge:
		metric, err := h.metSrv.GetGauge(ctx, metricName)

		if err != nil {
			// TODO: убрать костыль и сделать проверку на no rows
			log.Printf("unable to select gauge metric with name=%s", metricName)
			http.Error(w, fmt.Sprintf("unable to select gauge metric with name=%s", metricName),
				http.StatusNotFound)

			return
		}

		_, err = w.Write([]byte(strconv.FormatFloat(metric.Value, 'f', -1, 64)))

		if err != nil {
			log.Printf("unable to write body")
			http.Error(w, "internal server error",
				http.StatusNotFound)
		}
	case counter:
		metric, err := h.metSrv.GetCounter(ctx, metricName)

		if err != nil {
			// TODO: убрать костыль и сделать проверку на no rows
			log.Printf("unable to select counter metric with name=%s", metricName)
			http.Error(w, fmt.Sprintf("unable to counter gauge metric with name=%s", metricName),
				http.StatusNotFound)

			return
		}

		_, err = w.Write([]byte(fmt.Sprintf("%d", metric.Value)))

		if err != nil {
			log.Printf("unable to write body")
			http.Error(w, "internal server error",
				http.StatusNotFound)
		}
	default:
		log.Printf("unknown metric type. Expected %s or %s. Actual: %s\n", gauge, counter, metricType)
		http.Error(w, fmt.Sprintf("unknown metric type. Expected %s or %s. Actual: %s\n", gauge, counter, metricType),
			http.StatusNotImplemented)

		return
	}
}

// GetStaticAllMetrics return HTML with information about all metrics which exist
func (h *Handler) GetStaticAllMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tmpl := template.Must(template.New("metrics").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
</head>
<body>
	{{range . }}
		{{if .Name}}
<ol>{{.Name}}: {{.Value}}</ol>
		{{end}}
	{{end}}
</body>
</html>
`))

	data, err := h.metSrv.GetAll(ctx)

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	err = tmpl.Execute(w, data)

	if err != nil {
		log.Printf("template execute finished with err. Error: %s\n", err)

		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
