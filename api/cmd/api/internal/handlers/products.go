package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/ivorscott/go-delve-reload/internal/platform/database"
	"github.com/ivorscott/go-delve-reload/internal/platform/web"
	"github.com/ivorscott/go-delve-reload/internal/product"
	"github.com/pkg/errors"
)

// Products holds the application state needed by the handler methods.
type Products struct {
	repo *database.Repository
	log  *log.Logger
}

// List gets all products
func (p *Products) List(w http.ResponseWriter, r *http.Request) error {
	list, err := product.List(r.Context(), p.repo)
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, list, http.StatusOK)
}

// Retrieve a single Product
func (p *Products) Retrieve(w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")

	prod, err := product.Retrieve(r.Context(), p.repo, id)
	if err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for products %q", id)
		}
	}

	return web.Respond(r.Context(), w, prod, http.StatusOK)
}

// Create a new Product
func (p *Products) Create(w http.ResponseWriter, r *http.Request) error {

	var np product.NewProduct
	if err := web.Decode(r, &np); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	prod, err := product.Create(r.Context(), p.repo, np, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, prod, http.StatusCreated)
}

// Update decodes the body of a request to update an existing product. The ID
// of the product is part of the request URL.
func (p *Products) Update(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	var update product.UpdateProduct
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding product update")
	}

	if err := product.Update(r.Context(), p.repo, id, update, time.Now()); err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating product %q", id)
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusNoContent)
}

// Delete removes a single product identified by an ID in the request URL.
func (p *Products) Delete(w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")

	if err := product.Delete(r.Context(), p.repo, id); err != nil {
		switch err {
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "deleting product %q", id)
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusNoContent)
}
