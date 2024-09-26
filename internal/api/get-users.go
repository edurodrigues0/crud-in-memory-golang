package api

import (
	"crud-in-memory-golang/models"
	"net/http"
)

func handleGetUsers(db map[models.ID]models.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := make([]models.User, 0, len(db))
		for _, item := range db {
			users = append(users, item)
		}

		sendJSON(w, apiResponse{Data: users}, http.StatusOK)
	}
}
