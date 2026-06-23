package models

import "time"

type Cliente struct {
	ID        string    `json:"id"`
	Nombre    string    `json:"nombre"`
	Telefono  string    `json:"telefono"`
	CreatedAt time.Time `json:"created_at"`
}
