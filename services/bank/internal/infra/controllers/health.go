package controllers

import (
	"encoding/json"
	"net/http"
)

type healthHttpHandler struct {
	version string
}

func NewHealthHttpHandler(version string) *healthHttpHandler {
	return &healthHttpHandler{
		version: version,
	}
}

func (c *healthHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server-Version", c.version)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	result, _ := json.Marshal(map[string]string{
		"api": "OK",
	})
	w.Write(result)
}
