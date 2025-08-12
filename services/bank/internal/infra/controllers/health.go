package controllers

import "net/http"

func (controller *controller) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server-Version", controller.version)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
