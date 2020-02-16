package product

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/ivorscott/go-delve-reload/internal/platform/database"
	"github.com/pkg/errors"
)

// The product package shouldn't know anything about http
// While it may identify common know errors, how to respond is left to the handlers
var (
	ErrNotFound  = errors.New("product not found")
	ErrInvalidID = errors.New("id provided was not a valid UUID")
)

// List gets all Products from the database.
func List(ctx context.Context, repo *database.Repository) ([]Product, error) {
	products := []Product{}

	stmt := repo.SQ.Select(
		"id",
		"name",
		"price",
		"description",
		"created",
		"tags",
	).From(
		"products",
	)

	q, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "building query: %v", args)
	}

	if err := repo.DB.SelectContext(ctx, &products, q); err != nil {
		return nil, errors.Wrap(err, "selecting products")
	}

	return products, nil
}

// Retrieve finds the product identified by a given ID.
func Retrieve(ctx context.Context, repo *database.Repository, id string) (*Product, error) {
	var p Product

	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}

	stmt := repo.SQ.Select(
		"id",
		"name",
		"price",
		"description",
		"created",
		"tags",
	).From(
		"products",
	).Where(sq.Eq{"id": "?"})

	q, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "building query: %v", args)
	}

	if err := repo.DB.GetContext(ctx, &p, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &p, nil
}

// Create adds a new Product
func Create(ctx context.Context, repo *database.Repository, np NewProduct, now time.Time) (*Product, error) {
	p := Product{
		ID:          uuid.New().String(),
		Name:        np.Name,
		Price:       np.Price,
		Description: np.Description,
		Created:     now,
		Tags:        np.Tags,
	}

	stmt := repo.SQ.Insert(
		"products",
	).SetMap(map[string]interface{}{
		"id":          p.ID,
		"name":        p.Name,
		"price":       p.Price,
		"description": p.Description,
		"created":     p.Created,
		"tags":        p.Tags,
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return nil, errors.Wrapf(err, "inserting product: %v", np)
	}

	return &p, nil
}

// Update modifies data about a Product. It will error if the specified ID is
// invalid or does not reference an existing Product.
func Update(ctx context.Context, repo *database.Repository, id string, update UpdateProduct, now time.Time) error {
	p, err := Retrieve(ctx, repo, id)
	if err != nil {
		return err
	}

	if update.Name != nil {
		p.Name = *update.Name
	}
	if update.Price != nil {
		p.Price = *update.Price
	}
	if update.Description != nil {
		p.Description = *update.Description
	}
	if update.Tags != nil {
		p.Tags = update.Tags
	}

	stmt := repo.SQ.Update(
		"products",
	).SetMap(map[string]interface{}{
		"name":        p.Name,
		"price":       p.Price,
		"description": p.Description,
		"tags":        p.Tags,
	}).Where(sq.Eq{"id": id})

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return errors.Wrap(err, "updating product")
	}

	return nil
}

// Delete removes the product identified by a given ID.
func Delete(ctx context.Context, repo *database.Repository, id string) error {

	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidID
	}

	stmt := repo.SQ.Delete(
		"products",
	).Where(sq.Eq{"id": id})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return errors.Wrapf(err, "deleting product %s", id)
	}

	return nil
}
