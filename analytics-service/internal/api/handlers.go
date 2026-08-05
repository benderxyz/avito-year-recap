package api

import (
	"encoding/json"
	"net/http"

	"analytics-service/internal/aggregation"
)

func GetMetrics(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	metrics, ok := aggregation.GetByUserID(id)
	if !ok {
		http.Error(w, "metrics not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}
