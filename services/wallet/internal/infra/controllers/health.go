package controllers

import (
	"encoding/json"
	"net/http"
)

func (handler *APIHandler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server-Version", handler.version)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	result, _ := json.Marshal(map[string]string{
		"api": "OK",
	})
	w.Write(result)
}
