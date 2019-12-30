package models

import (
	"time"
)

// Products represents product data from the database
type Product struct {
	ID          int
	Name        string
	Price       int
	Description string
	Created     time.Time
}
