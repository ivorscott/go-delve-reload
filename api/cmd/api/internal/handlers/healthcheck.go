package handlers

import (
	"net/http"

	"github.com/ivorscott/go-delve-reload/internal/platform/database"
	"github.com/ivorscott/go-delve-reload/internal/platform/web"
)

// Healthcheck supports orchestration
type HealthCheck struct {
	repo *database.Repository
}

// Health validates the service is healthy and ready to accept requests.
func (c *HealthCheck) Health(w http.ResponseWriter, r *http.Request) error {

	var health struct {
		Status string `json:"status"`
	}

	// Check if the database is ready.
	if err := database.StatusCheck(r.Context(), c.repo.DB); err != nil {

		// If the database is not ready we will tell the client and use a 500
		// status. Do not respond by just returning an error because further up in
		// the call stack will interpret that as an unhandled error.
		health.Status = "db not ready"
		return web.Respond(r.Context(), w, health, http.StatusInternalServerError)
	}

	health.Status = "ok"
	return web.Respond(r.Context(), w, health, http.StatusOK)
}
