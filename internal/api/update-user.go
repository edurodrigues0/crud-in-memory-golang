package api

import (
	"crud-in-memory-golang/models"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PutBody struct {
	FirstName string `json:"firstname" validate:"required"`
	LastName  string `json:"lastname" validate:"required"`
	Biography string `json:"biography" validate:"required"`
}

func handleUpdateUser(db map[models.ID]models.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, _ := uuid.Parse(idStr)

		var body PutBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendJSON(
				w,
				apiResponse{Error: "Invalid body"},
				http.StatusUnprocessableEntity,
			)
			return
		}

		if err := validate.Struct(body); err != nil {
			sendJSON(
				w,
				apiResponse{Error: "Missing required fields"},
				http.StatusBadRequest,
			)
			return
		}

		userID := models.ID(id)
		_, ok := db[userID]
		if ok {
			db[userID] = models.User{
				FirstName: body.FirstName,
				LastName:  body.LastName,
				Biography: body.Biography,
			}

			userUpdated := db[userID]

			sendJSON(
				w,
				apiResponse{Data: userUpdated},
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
