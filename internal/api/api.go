package api

import (
	"crud-in-memory-golang/models"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type apiResponse struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

func NewHandler(db map[models.ID]models.User) http.Handler {
	r := chi.NewMux()

	r.Use(
		middleware.Recoverer,
		middleware.RequestID,
		middleware.Logger,
	)

	r.Route("/api/users", func(r chi.Router) {
		r.Post("/", handleCreateUser(db))
		r.Get("/", handleGetUsers(db))

		r.Route("/{id}", func(r chi.Router) {
			r.Use(ValidateUUID)
			r.Get("/", handleGetUser(db))
			r.Delete("/", handleDeleteUser(db))
			r.Put("/", handleUpdateUser(db))
		})
	})

	return r
}
