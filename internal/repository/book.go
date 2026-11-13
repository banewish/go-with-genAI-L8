package repository

import (
	"context"

	"libratrack/internal/model"
)

type BookRepository interface {
	Create(ctx context.Context, req *model.BookRequest) (*model.Book, error)
	FindByID(ctx context.Context, id uint) (*model.Book, error)
	List(ctx context.Context, offset, limit int) ([]model.Book, error)
	Update(ctx context.Context, id uint, req *model.BookRequest) (*model.Book, error)
	Delete(ctx context.Context, id uint) error
}
