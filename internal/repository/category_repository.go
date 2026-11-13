package repository

import "context"
import "libratrack/internal/model"

type CategoryRepository interface {
	List(ctx context.Context, offset, limit int) ([]model.Category, error)
	FindByID(ctx context.Context, id uint) (*model.Category, error)
	Create(ctx context.Context, req *model.Category) (*model.Category, error)
	Update(ctx context.Context, id uint, req *model.Category) (*model.Category, error)
	Delete(ctx context.Context, id uint) error
}
