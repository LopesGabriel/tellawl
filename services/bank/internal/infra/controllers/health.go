package controllers

import "net/http"

type healthHttpHandler struct {
	version string
}

func NewHealthHttpHandler() *healthHttpHandler {
	return &healthHttpHandler{
		version: "1.0.0",
	}
}

func (c *healthHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server-Version", c.version)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
