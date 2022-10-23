package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

// Handler implementing all handlers for server
type Handler struct {
}

// Config for Handler
type Config struct {
	Router *mux.Router
}

// New is constructor for Handler
func New(c *Config) {
	h := Handler{}

	c.Router.HandleFunc("/update/{metric_type}/{metric_name}/{value}", h.UpdateMetric).Methods(http.MethodPost)
}

// UpdateMetric accepts request for create or update metrics
func (h *Handler) UpdateMetric(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	metricType, ok := vars["metric_type"]

	if !ok {
		http.Error(w, "unable parse path parameter 'metric_type'", http.StatusBadRequest)
	}

	metricName, ok := vars["metric_name"]

	if !ok {
		http.Error(w, "unable parse path parameter 'metric_name'", http.StatusBadRequest)
	}

	value, ok := vars["value"]

	if !ok {
		http.Error(w, "unable parse path parameter 'value'", http.StatusBadRequest)
	}

	fmt.Println(metricType, metricName, value)

	_, err := w.Write([]byte("OK"))

	if err != nil {
		w.WriteHeader(500)

		return
	}
}
