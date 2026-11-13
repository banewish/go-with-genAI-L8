package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"homework/internal/model"
	"homework/internal/repository"
)

type BookHandler struct {
	repo *http.ServeMux
}

type bookHandlerImpl struct {
	repo repository.BookRepository
}

func NewBookHandler(repo repository.BookRepository) http.Handler {
	impl := &bookHandlerImpl{repo: repo}
	mux := http.NewServeMux()
	log.Printf("Registering book handlers...")
	mux.HandleFunc("/api/v1/books", impl.handleBooks)
	mux.HandleFunc("/api/v1/books/", impl.handleBooks)
	return mux
}

func (h *bookHandlerImpl) handleBooks(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/books")

	switch {
	case r.Method == http.MethodGet && (path == "" || path == "/"):
		h.ListBooks(w, r)
	case r.Method == http.MethodGet && path != "" && path != "/":
		h.GetBook(w, r)
	case r.Method == http.MethodPost && (path == "" || path == "/"):
		h.CreateBook(w, r)
	case r.Method == http.MethodPut && path != "" && path != "/":
		h.UpdateBook(w, r)
	case r.Method == http.MethodDelete && path != "" && path != "/":
		h.DeleteBook(w, r)
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func (h *bookHandlerImpl) ListBooks(w http.ResponseWriter, r *http.Request) {
	offset := 0
	limit := 10

	if q := r.URL.Query().Get("offset"); q != "" {
		if v, err := strconv.Atoi(q); err == nil && v >= 0 {
			offset = v
		}
	}
	if q := r.URL.Query().Get("limit"); q != "" {
		if v, err := strconv.Atoi(q); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}

	books, err := h.repo.List(r.Context(), offset, limit)
	if err != nil {
		log.Printf("Error listing books: %v", err)
		http.Error(w, "Failed to list books", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": books,
		"meta": map[string]int{
			"offset": offset,
			"limit":  limit,
			"count":  len(books),
		},
	})
}

func (h *bookHandlerImpl) GetBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(strings.TrimPrefix(r.URL.Path, "/api/v1/books/"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	book, err := h.repo.FindByID(r.Context(), uint(id))
	if err != nil {
		log.Printf("Error finding book %d: %v", id, err)
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *bookHandlerImpl) CreateBook(w http.ResponseWriter, r *http.Request) {
	var req model.BookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	book, err := h.repo.Create(r.Context(), &req)
	if err != nil {
		log.Printf("Error creating book: %v", err)
		http.Error(w, "Failed to create book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func (h *bookHandlerImpl) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(strings.TrimPrefix(r.URL.Path, "/api/v1/books/"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var req model.BookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	book, err := h.repo.Update(r.Context(), uint(id), &req)
	if err != nil {
		log.Printf("Error updating book %d: %v", id, err)
		http.Error(w, "Failed to update book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *bookHandlerImpl) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(strings.TrimPrefix(r.URL.Path, "/api/v1/books/"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(r.Context(), uint(id)); err != nil {
		log.Printf("Error deleting book %d: %v", id, err)
		http.Error(w, "Failed to delete book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
