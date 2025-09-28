package api

import (
	"encoding/json"
	"net/http"

	"go.opentelemetry.io/otel/codes"
)

func (handler *apiHandler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	_, span := handler.tracer.Start(r.Context(), "HandleHealthCheck")
	defer span.End()

	result, _ := json.Marshal(map[string]string{
		"api": "OK",
	})

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
	span.SetStatus(codes.Ok, "OK")
}
