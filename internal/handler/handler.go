package handler

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mtrrun/internal/model"
)

const (
	metricTypeGauge   = "gauge"
	metricTypeCounter = "counter"
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

	c.Router.HandleFunc("/", panicMiddleware(h.GetStaticAllMetrics)).Methods(http.MethodGet)
	c.Router.HandleFunc("/update/{metric_type}/{metric_name}/{value}", panicMiddleware(h.UpdateMetric)).Methods(http.MethodPost)
	c.Router.HandleFunc("/value/{metric_type}/{metric_name}", panicMiddleware(h.GetMetric)).Methods(http.MethodGet)
}

// For recover in request process with panic
func panicMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if i := recover(); i != nil {
				log.Printf("panic at %s: %v\n", r.URL.Path, i)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// UpdateMetric accepts request for create or update metrics
func (h *Handler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	metricType, ok := vars["metric_type"]

	if !ok {
		http.Error(w, "unable parse path parameter 'metric_type': parameter not found in request", http.StatusBadRequest)

		return
	}

	metricName, ok := vars["metric_name"]

	if !ok {
		http.Error(w, "unable parse path parameter 'metric_name': parameter not found in request", http.StatusBadRequest)

		return
	}

	value, ok := vars["value"]

	if !ok {
		http.Error(w, "unable parse path parameter 'value': parameter not found in request", http.StatusBadRequest)

		return
	}

	// Validate path params

	if len(metricName) == 0 {
		msg := "unable to parse name. Expected: string with length > 0"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)

		return
	}

	switch metricType {
	case metricTypeGauge:
		valueFloat64, err := strconv.ParseFloat(value, 64)

		if err != nil {
			msg := fmt.Sprintf("unable to parse value. Expected: float. Actual: %s", value)
			log.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)

			return
		}

		err = h.metSrv.PutGauge(ctx, model.PutGaugeDTO{
			Name:  metricName,
			Value: valueFloat64,
		})

		if err != nil {
			msg := fmt.Sprintf("unable to update/create gauge metric with name=%s and value=%s", metricName, value)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)

			return
		}
	case metricTypeCounter:
		valueInt64, err := strconv.ParseInt(value, 10, 64)

		if err != nil {
			msg := fmt.Sprintf("unable to parse value. Expected: int. Actual: %s", value)
			log.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)

			return
		}

		err = h.metSrv.PutCounter(ctx, model.PutCounterDTO{
			Name:  metricName,
			Value: valueInt64,
		})

		if err != nil {
			msg := fmt.Sprintf("unable to update/create counter metric with name=%s and value=%s", metricName, value)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)

			return
		}
	default:
		msg := fmt.Sprintf("unknown metric type. Expected %s or %s. Actual: %s\n", metricTypeGauge, metricTypeCounter, metricType)
		log.Println(msg)
		http.Error(w, msg, http.StatusNotImplemented)

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
		http.Error(w, "unable parse path parameter 'metric_type': parameter not found in request", http.StatusBadRequest)

		return
	}

	metricName, ok := vars["metric_name"]

	if !ok {
		http.Error(w, "unable parse path parameter 'metric_name': parameter not found in request", http.StatusBadRequest)

		return
	}

	if len(metricName) == 0 {
		http.Error(w, "unable parse path parameter 'metric_name'. Expected: string with length > 0", http.StatusBadRequest)

		return
	}

	switch metricType {
	case metricTypeGauge:
		metric, err := h.metSrv.GetGauge(ctx, metricName)

		if err != nil {
			// TODO: убрать костыль и сделать проверку на no rows
			msg := fmt.Sprintf("unable to select gauge metric with name=%s", metricName)
			log.Println(msg)
			http.Error(w, msg, http.StatusNotFound)

			return
		}

		_, err = w.Write([]byte(strconv.FormatFloat(metric.Value, 'f', -1, 64)))

		if err != nil {
			log.Printf("unable to write body. Error: %s\n", err)
			http.Error(w, "internal server error",
				http.StatusNotFound)
		}
	case metricTypeCounter:
		metric, err := h.metSrv.GetCounter(ctx, metricName)

		if err != nil {
			// TODO: убрать костыль и сделать проверку на no rows
			msg := fmt.Sprintf("unable to select counter metric with name=%s", metricName)
			log.Println(msg)
			http.Error(w, msg, http.StatusNotFound)

			return
		}

		_, err = w.Write([]byte(fmt.Sprintf("%d", metric.Value)))

		if err != nil {
			log.Printf("unable to write body. Error: %s\n", err)
			http.Error(w, "internal server error",
				http.StatusNotFound)
		}
	default:
		msg := fmt.Sprintf("unknown metric type. Expected %s or %s. Actual: %s", metricTypeGauge, metricTypeCounter, metricType)
		log.Println(msg)
		http.Error(w, msg, http.StatusNotImplemented)

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
		log.Printf("unable to get all metrics. Error: %s\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	err = tmpl.Execute(w, data)

	if err != nil {
		log.Printf("template execute finished with err. Error: %s\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
