package controllers

import (
	"encoding/json"
	"net/http"
)

func WriteError(w http.ResponseWriter, statusCode int, payload map[string]any) {
	result, _ := json.Marshal(payload)

	w.Header().Del("Content-Type")
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(result)
}
