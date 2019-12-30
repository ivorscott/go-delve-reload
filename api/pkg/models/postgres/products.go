package postgres

import (
	"github.com/ivorscott/go-delve-reload/pkg/models"
	"github.com/jinzhu/gorm"
)

// ProductModel object
type ProductModel struct {
	DB *gorm.DB
}

// Get all products
func (m *ProductModel) GetAll() []*models.Product {
	products := []*models.Product{}
	m.DB.Find(&products)
	return products
}
