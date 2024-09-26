package api

import (
	"crud-in-memory-golang/models"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func handleDeleteUser(db map[models.ID]models.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, _ := uuid.Parse(idStr)

		userID := models.ID(id)
		_, ok := db[userID]
		if !ok {
			sendJSON(
				w,
				apiResponse{Error: "The user could not be removed"},
				http.StatusNotFound,
			)
			return
		}

		delete(db, userID)
	}
}
