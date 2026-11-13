package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"libratrack/internal/model"
	"libratrack/internal/repository"
)

type PetHandler struct {
	repo *repository.PetRepository
}

func NewPetHandler(repo *repository.PetRepository) *PetHandler {
	return &PetHandler{repo: repo}
}

func (h *PetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	path := r.URL.Path
	method := r.Method

	if path == "/health" {
		if method == http.MethodGet {
			h.health(w)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if !strings.HasPrefix(path, "/api/v1/pets") {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if path == "/api/v1/pets" {
		switch method {
		case http.MethodPost:
			h.create(w, r)
		case http.MethodGet:
			h.list(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	idPart := strings.TrimPrefix(path, "/api/v1/pets/")
	if idPart == path || idPart == "" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if strings.Contains(idPart, "/") {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	id, err := strconv.ParseUint(idPart, 10, 32)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid_id", "Invalid ID")
		return
	}

	if id > 4294967295 {
		jsonError(w, http.StatusBadRequest, "invalid_id", "Invalid ID")
		return
	}

	switch method {
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

func (h *PetHandler) health(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *PetHandler) create(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		jsonError(w, http.StatusUnprocessableEntity, "validation_error", "Invalid JSON")
		return
	}

	if len(body) == 0 {
		jsonError(w, http.StatusUnprocessableEntity, "validation_error", "Invalid JSON")
		return
	}

	var req model.PetRequest
	if err := json.Unmarshal(body, &req); err != nil {
		jsonError(w, http.StatusUnprocessableEntity, "validation_error", "Invalid JSON")
		return
	}

	// Check if body was an array
	if body[0] == '[' {
		jsonError(w, http.StatusUnprocessableEntity, "validation_error", "Invalid JSON")
		return
	}

	if err := validatePetRequest(&req); err != "" {
		jsonError(w, http.StatusUnprocessableEntity, "validation_error", err)
		return
	}

	pet := h.repo.Create(&req)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pet)
}

func (h *PetHandler) get(w http.ResponseWriter, id uint32) {
	pet := h.repo.GetByID(id)
	if pet == nil {
		jsonError(w, http.StatusNotFound, "not_found", "Pet not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pet)
}

func (h *PetHandler) list(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	speciesFilter := r.URL.Query().Get("species")

	page := 1
	limit := 10

	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			jsonError(w, http.StatusBadRequest, "invalid_pagination", "Invalid pagination")
			return
		}
		page = p
	}

	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l < 1 {
			jsonError(w, http.StatusBadRequest, "invalid_pagination", "Invalid pagination")
			return
		}
		if l > 100 {
			l = 100
		}
		limit = l
	}

	offset := (page - 1) * limit
	pets := h.repo.List(offset, limit, speciesFilter)

	if pets == nil {
		pets = []*model.Pet{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pets)
}

func (h *PetHandler) update(w http.ResponseWriter, r *http.Request, id uint32) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid_json", "Invalid JSON")
		return
	}

	if len(body) == 0 {
		jsonError(w, http.StatusUnprocessableEntity, "validation_error", "Invalid JSON")
		return
	}

	var req model.PetRequest
	if err := json.Unmarshal(body, &req); err != nil {
		jsonError(w, http.StatusUnprocessableEntity, "validation_error", "Invalid JSON")
		return
	}

	// Check if body was an array
	if body[0] == '[' {
		jsonError(w, http.StatusUnprocessableEntity, "validation_error", "Invalid JSON")
		return
	}

	if err := validatePetRequest(&req); err != "" {
		jsonError(w, http.StatusUnprocessableEntity, "validation_error", err)
		return
	}

	pet := h.repo.Update(id, &req)
	if pet == nil {
		jsonError(w, http.StatusNotFound, "not_found", "Pet not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pet)
}

func (h *PetHandler) delete(w http.ResponseWriter, id uint32) {
	if !h.repo.Delete(id) {
		jsonError(w, http.StatusNotFound, "not_found", "Pet not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func validatePetRequest(req *model.PetRequest) string {
	if err := validateRequired(req.Name); err != "" {
		return err
	}
	if err := validateRequired(req.Owner); err != "" {
		return err
	}
	if err := validateRequired(req.Breed); err != "" {
		return err
	}
	if err := validateRequired(req.Species); err != "" {
		return err
	}

	if req.Age != nil && *req.Age < 0 {
		return "age must be non-negative"
	}

	return ""
}

func validateRequired(s string) string {
	if s == "" || strings.TrimSpace(s) == "" {
		return "required field is empty"
	}
	return ""
}

func jsonError(w http.ResponseWriter, status int, code, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]map[string]string{
		"error": {
			"code":    code,
			"message": msg,
		},
	})
}
