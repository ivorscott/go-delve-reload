package product

import (
	"time"
)

// Products represents product data from the database
type Product struct {
	ID          string    `db:"id" json:"id" `
	Name        string    `db:"name" json:"name"`
	Price       int       `db:"price" json:"price"`
	Description string    `db:"description" json:"description"`
	Created     time.Time `db:"created" json:"created"`
	Tags        *string   `db:"tags" json:"tags"`
}

type NewProduct struct {
	Name        string  `json:"name" validate:"required"`
	Price       int     `json:"price" validate:"gte=0"`
	Description string  `json:"description"`
	Tags        *string `json:"tags"`
}

type UpdateProduct struct {
	Name        *string `json:"name"`
	Price       *int    `json:"price" validate:"omitempty,gte=0"`
	Description *string `json:"description"`
	Tags        *string `json:"tags"`
}
