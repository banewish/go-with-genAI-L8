package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"homework/internal/handler"
	"homework/internal/model"
	"homework/internal/repository"
	"context"
	"gorm.io/gorm"
)

type mockHandlerRepo struct {
	books []model.Book
}

func (m *mockHandlerRepo) List(ctx context.Context, offset, limit int) ([]model.Book, error) {
	return m.books, nil
}

func (m *mockHandlerRepo) FindByID(ctx context.Context, id uint) (*model.Book, error) {
	for _, b := range m.books {
		if b.ID == id {
			return &b, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockHandlerRepo) Create(ctx context.Context, req *model.BookRequest) (*model.Book, error) {
	return &model.Book{}, nil
}

func (m *mockHandlerRepo) Update(ctx context.Context, id uint, req *model.BookRequest) (*model.Book, error) {
	return &model.Book{}, nil
}

func (m *mockHandlerRepo) Delete(ctx context.Context, id uint) error {
	return nil
}

func TestListBooks(t *testing.T) {
	repo := &mockHandlerRepo{
		books: []model.Book{{ID: 1, Title: "Test"}},
	}

	h := handler.NewBookHandler(repo)
	req := httptest.NewRequest("GET", "/api/v1/books", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
