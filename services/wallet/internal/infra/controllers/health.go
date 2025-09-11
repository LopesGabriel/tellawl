package controllers

import (
	"encoding/json"
	"net/http"
)

func (handler *APIHandler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "HandleHealthCheck")
	defer span.End()

	result, _ := json.Marshal(map[string]string{
		"api": "OK",
	})

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
