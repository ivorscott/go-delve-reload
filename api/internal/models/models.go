package models

import (
	"time"
)

// Products represents product data from the database
type Product struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Price       int       `json:"price"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
}
