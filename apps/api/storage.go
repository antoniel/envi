package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUser(user *User) error
	GetUser(id string) (*User, error)
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type PostgresStorage struct {
	db *sql.DB
}
