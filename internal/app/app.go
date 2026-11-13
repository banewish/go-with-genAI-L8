// Package app wires the whole HTTP API together.
//
// The automatic tests talk to your API ONLY through NewRouter(), so you are free to
// organise the rest of the code (model, repository, handler, middleware) as you like.
package app

import (
	"net/http"

	"gorm.io/gorm"
	"homework/internal/handler"
	"homework/internal/repository"
)

// NewRouter must return a fully configured HTTP handler for your resource.
//
// Requirements (see README.md for the full contract):
//   - every call returns a NEW router with its own EMPTY in-memory storage;
//   - it must be safe for concurrent requests;
//   - it can be built with net/http, gin, chi or any other router.
//
// TODO: replace this stub with your implementation.
func NewRouter() http.Handler {
	repo := repository.NewPetRepository()
	mux := http.NewServeMux()
	books := handler.NewLegacyBookHandler()
	mux.Handle("/api/v1/books", books)
	mux.Handle("/api/v1/books/", books)
	mux.Handle("/", handler.NewPetHandler(repo))
	return mux
}

// NewRouterWithDB creates a router with PostgreSQL database connection
func NewRouterWithDB(db *gorm.DB) http.Handler {
	bookRepo := repository.NewGormBookRepository(db)
	bookHandler := handler.NewBookHandler(bookRepo)
	return bookHandler
}
