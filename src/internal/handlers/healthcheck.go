package handlers

import "net/http"

type HealthCheckHandler struct{}

func NewHealthcheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{}
}

func (h *HealthCheckHandler) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))

}
