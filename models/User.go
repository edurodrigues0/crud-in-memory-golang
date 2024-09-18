package models

import "github.com/google/uuid"

type ID uuid.UUID

type User struct {
	FirstName string
	LastName  string
	Biography string
}

type Application struct {
	data map[ID]User
}
