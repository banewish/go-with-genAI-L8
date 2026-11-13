package repository

import "context"
import "libratrack/internal/model"

type AuthorRepository interface {
	List(ctx context.Context, offset, limit int) ([]model.Author, error)
	FindByID(ctx context.Context, id uint) (*model.Author, error)
	Create(ctx context.Context, req *model.Author) (*model.Author, error)
	Update(ctx context.Context, id uint, req *model.Author) (*model.Author, error)
	Delete(ctx context.Context, id uint) error
}
