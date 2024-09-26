package api

import (
	"crud-in-memory-golang/models"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New()

type PostBody struct {
	FirstName string `json:"firstname" validate:"required"`
	LastName  string `json:"lastname" validate:"required"`
	Biography string `json:"biography" validate:"required"`
}

func handleCreateUser(db map[models.ID]models.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body PostBody
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

		ID, err := uuid.NewUUID()
		if err != nil {
			sendJSON(
				w,
				apiResponse{Error: "Error on create uuid"},
				http.StatusBadRequest,
			)
			return
		}

		userID := models.ID(ID)

		if _, ok := db[userID]; ok {
			sendJSON(
				w,
				apiResponse{Error: "ID already exists"},
				http.StatusConflict,
			)
			return
		}

		db[userID] = models.User{
			FirstName: body.FirstName,
			LastName:  body.LastName,
			Biography: body.Biography,
		}

		sendJSON(
			w,
			apiResponse{Data: ID},
			http.StatusCreated,
		)
		defer r.Body.Close()
	}
}
