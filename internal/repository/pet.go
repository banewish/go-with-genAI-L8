package repository

import (
	"sync"
	"time"

	"homework/internal/model"
)

type PetRepository struct {
	mu     sync.RWMutex
	pets   map[uint32]*model.Pet
	nextID uint32
}

func NewPetRepository() *PetRepository {
	return &PetRepository{
		pets:   make(map[uint32]*model.Pet),
		nextID: 1,
	}
}

func (r *PetRepository) Create(req *model.PetRequest) *model.Pet {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()
	pet := &model.Pet{
		ID:        r.nextID,
		Name:      req.Name,
		Owner:     req.Owner,
		Breed:     req.Breed,
		Species:   req.Species,
		Age:       req.Age,
		Notes:     req.Notes,
		CreatedAt: now,
		UpdatedAt: now,
	}

	r.pets[r.nextID] = pet
	r.nextID++
	return pet
}

func (r *PetRepository) GetByID(id uint32) *model.Pet {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.pets[id]
}

func (r *PetRepository) List(offset, limit int, speciesFilter string) []*model.Pet {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*model.Pet
	for _, pet := range r.pets {
		if speciesFilter == "" || pet.Species == speciesFilter {
			// Copy the pet to avoid race conditions
			petCopy := *pet
			results = append(results, &petCopy)
		}
	}

	// Sort by ID
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].ID < results[i].ID {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	start := offset
	end := offset + limit
	if start > len(results) {
		return []*model.Pet{}
	}
	if end > len(results) {
		end = len(results)
	}
	return results[start:end]
}

func (r *PetRepository) Update(id uint32, req *model.PetRequest) *model.Pet {
	r.mu.Lock()
	defer r.mu.Unlock()

	pet, exists := r.pets[id]
	if !exists {
		return nil
	}

	pet.Name = req.Name
	pet.Owner = req.Owner
	pet.Breed = req.Breed
	pet.Species = req.Species
	pet.Age = req.Age
	pet.Notes = req.Notes
	pet.UpdatedAt = time.Now().UTC()

	return pet
}

func (r *PetRepository) Delete(id uint32) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.pets[id]; !exists {
		return false
	}
	delete(r.pets, id)
	return true
}

func (r *PetRepository) Count(speciesFilter string) int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, pet := range r.pets {
		if speciesFilter == "" || pet.Species == speciesFilter {
			count++
		}
	}
	return count
}
