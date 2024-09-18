package api

import (
	"crud-in-memory-golang/models"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New()

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

type Response struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

func sendJSON(w http.ResponseWriter, resp Response, status int) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(resp)
	if err != nil {
		slog.Error("Failed to marshal json data", "error", err)
		sendJSON(
			w,
			Response{Error: "Something went wrong"},
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error("Failed to write response to client", "error", err)
		return
	}
}

func ValidateUUID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if _, err := uuid.Parse(id); err != nil {
			sendJSON(
				w,
				Response{Error: "Invalid UUID"},
				http.StatusBadRequest,
			)
			return
		}
		next.ServeHTTP(w, r)
	})
}

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
				Response{Error: "Invalid body"},
				http.StatusUnprocessableEntity,
			)
			return
		}

		if err := validate.Struct(body); err != nil {
			sendJSON(
				w,
				Response{Error: "Missing required fields"},
				http.StatusBadRequest,
			)
			return
		}

		ID, err := uuid.NewUUID()
		if err != nil {
			sendJSON(
				w,
				Response{Error: "Error on create uuid"},
				http.StatusBadRequest,
			)
			return
		}

		userID := models.ID(ID)

		if _, ok := db[userID]; ok {
			sendJSON(
				w,
				Response{Error: "ID already exists"},
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
			Response{Data: ID},
			http.StatusCreated,
		)
		defer r.Body.Close()
	}
}

func handleGetUsers(db map[models.ID]models.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := make([]models.User, 0, len(db))
		for _, item := range db {
			users = append(users, item)
		}

		sendJSON(w, Response{Data: users}, http.StatusOK)
	}
}

func handleGetUser(db map[models.ID]models.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, _ := uuid.Parse(idStr)

		userID := models.ID(id)
		user, ok := db[userID]
		if ok {
			sendJSON(
				w,
				Response{Data: user},
				http.StatusOK,
			)
			return
		}

		sendJSON(
			w,
			Response{Error: "User not found."},
			http.StatusNotFound,
		)
	}
}

func handleDeleteUser(db map[models.ID]models.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, _ := uuid.Parse(idStr)

		userID := models.ID(id)
		_, ok := db[userID]
		if !ok {
			sendJSON(
				w,
				Response{Error: "The user could not be removed"},
				http.StatusNotFound,
			)
			return
		}

		delete(db, userID)
	}
}

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
				Response{Error: "Invalid body"},
				http.StatusUnprocessableEntity,
			)
			return
		}

		if err := validate.Struct(body); err != nil {
			sendJSON(
				w,
				Response{Error: "Missing required fields"},
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
				Response{Data: userUpdated},
				http.StatusOK,
			)
			return
		}

		sendJSON(
			w,
			Response{Error: "User not found."},
			http.StatusNotFound,
		)
	}
}
