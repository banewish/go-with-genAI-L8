package repository_test

import (
	"context"
	"testing"
	"libratrack/internal/model"
	"libratrack/internal/repository"
	"gorm.io/gorm"
)

type mockBookRepository struct {
	books []model.Book
}

func (m *mockBookRepository) List(ctx context.Context, offset, limit int) ([]model.Book, error) {
	return m.books, nil
}

func (m *mockBookRepository) FindByID(ctx context.Context, id uint) (*model.Book, error) {
	for _, b := range m.books {
		if b.ID == id {
			return &b, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockBookRepository) Create(ctx context.Context, req *model.BookRequest) (*model.Book, error) {
	book := &model.Book{ID: 1, Title: req.Title}
	return book, nil
}

func (m *mockBookRepository) Update(ctx context.Context, id uint, req *model.BookRequest) (*model.Book, error) {
	book := &model.Book{ID: id, Title: req.Title}
	return book, nil
}

func (m *mockBookRepository) Delete(ctx context.Context, id uint) error {
	return nil
}

func TestRepository(t *testing.T) {
	repo := &mockBookRepository{
		books: []model.Book{{ID: 1, Title: "Test"}},
	}

	books, err := repo.List(context.Background(), 0, 10)
	if err != nil || len(books) == 0 {
		t.Error("Expected to find books")
	}
}
