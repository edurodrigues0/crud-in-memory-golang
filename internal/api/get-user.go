package api

import (
	"crud-in-memory-golang/models"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func handleGetUser(db map[models.ID]models.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, _ := uuid.Parse(idStr)

		userID := models.ID(id)
		user, ok := db[userID]
		if ok {
			sendJSON(
				w,
				apiResponse{Data: user},
				http.StatusOK,
			)
			return
		}

		sendJSON(
			w,
			apiResponse{Error: "User not found."},
			http.StatusNotFound,
		)
	}
}
