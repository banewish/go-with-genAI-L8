package handler

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type legacyBookRequest struct {
	Title         string  `json:"title"`
	ISBN          string  `json:"isbn"`
	Author        string  `json:"author"`
	Category      string  `json:"category"`
	PublishedYear *int    `json:"published_year"`
	Description   *string `json:"description"`
}

type legacyBook struct {
	legacyBookRequest
	ID        uint32    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type legacyBookHandler struct {
	mu     sync.RWMutex
	books  map[uint32]legacyBook
	nextID uint32
}

func NewLegacyBookHandler() http.Handler {
	return &legacyBookHandler{books: make(map[uint32]legacyBook), nextID: 1}
}

func (h *legacyBookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	const base = "/api/v1/books"
	if r.URL.Path == base {
		switch r.Method {
		case http.MethodGet:
			h.list(w, r)
		case http.MethodPost:
			h.create(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	idPart := strings.TrimPrefix(r.URL.Path, base+"/")
	if idPart == r.URL.Path || idPart == "" || strings.Contains(idPart, "/") {
		http.NotFound(w, r)
		return
	}
	id, err := strconv.ParseUint(idPart, 10, 32)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid_id", "Invalid ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.get(w, uint32(id))
	case http.MethodPut:
		h.update(w, r, uint32(id))
	case http.MethodDelete:
		h.delete(w, uint32(id))
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func decodeLegacyBook(w http.ResponseWriter, r *http.Request) (legacyBookRequest, bool) {
	var req legacyBookRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		jsonError(w, http.StatusUnprocessableEntity, "validation_error", "Invalid JSON")
		return req, false
	}
	for _, field := range []string{req.Title, req.ISBN, req.Author, req.Category} {
		if strings.TrimSpace(field) == "" {
			jsonError(w, http.StatusUnprocessableEntity, "validation_error", "Required field is empty")
			return req, false
		}
	}
	if req.PublishedYear != nil && *req.PublishedYear < 0 {
		jsonError(w, http.StatusUnprocessableEntity, "validation_error", "Published year must be non-negative")
		return req, false
	}
	return req, true
}

func writeLegacyBook(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(value)
}

func (h *legacyBookHandler) create(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeLegacyBook(w, r)
	if !ok {
		return
	}
	h.mu.Lock()
	now := time.Now().UTC()
	book := legacyBook{legacyBookRequest: req, ID: h.nextID, CreatedAt: now, UpdatedAt: now}
	h.books[book.ID] = book
	h.nextID++
	h.mu.Unlock()
	writeLegacyBook(w, http.StatusCreated, book)
}

func (h *legacyBookHandler) get(w http.ResponseWriter, id uint32) {
	h.mu.RLock()
	book, ok := h.books[id]
	h.mu.RUnlock()
	if !ok {
		jsonError(w, http.StatusNotFound, "not_found", "Book not found")
		return
	}
	writeLegacyBook(w, http.StatusOK, book)
}

func (h *legacyBookHandler) list(w http.ResponseWriter, r *http.Request) {
	page, limit := 1, 10
	for name, target := range map[string]*int{"page": &page, "limit": &limit} {
		if raw := r.URL.Query().Get(name); raw != "" {
			value, err := strconv.Atoi(raw)
			if err != nil || value < 1 {
				jsonError(w, http.StatusBadRequest, "invalid_pagination", "Invalid pagination")
				return
			}
			*target = value
		}
	}
	if limit > 100 {
		limit = 100
	}
	category := r.URL.Query().Get("category")
	books := make([]legacyBook, 0)
	h.mu.RLock()
	for _, book := range h.books {
		if category == "" || book.Category == category {
			books = append(books, book)
		}
	}
	h.mu.RUnlock()
	sort.Slice(books, func(i, j int) bool { return books[i].ID < books[j].ID })
	if page > len(books)/limit+1 {
		writeLegacyBook(w, http.StatusOK, []legacyBook{})
		return
	}
	start := (page - 1) * limit
	if start >= len(books) {
		writeLegacyBook(w, http.StatusOK, []legacyBook{})
		return
	}
	end := start + limit
	if end > len(books) {
		end = len(books)
	}
	writeLegacyBook(w, http.StatusOK, books[start:end])
}

func (h *legacyBookHandler) update(w http.ResponseWriter, r *http.Request, id uint32) {
	req, ok := decodeLegacyBook(w, r)
	if !ok {
		return
	}
	h.mu.Lock()
	book, found := h.books[id]
	if found {
		book.legacyBookRequest = req
		book.UpdatedAt = time.Now().UTC()
		h.books[id] = book
	}
	h.mu.Unlock()
	if !found {
		jsonError(w, http.StatusNotFound, "not_found", "Book not found")
		return
	}
	writeLegacyBook(w, http.StatusOK, book)
}

func (h *legacyBookHandler) delete(w http.ResponseWriter, id uint32) {
	h.mu.Lock()
	_, found := h.books[id]
	delete(h.books, id)
	h.mu.Unlock()
	if !found {
		jsonError(w, http.StatusNotFound, "not_found", "Book not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
