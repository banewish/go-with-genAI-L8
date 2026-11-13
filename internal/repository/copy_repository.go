package repository

import "context"
import "libratrack/internal/model"

type CopyRepository interface {
	List(ctx context.Context, offset, limit int) ([]model.Copy, error)
	FindByID(ctx context.Context, id uint) (*model.Copy, error)
	Create(ctx context.Context, req *model.Copy) (*model.Copy, error)
	Update(ctx context.Context, id uint, req *model.Copy) (*model.Copy, error)
	Delete(ctx context.Context, id uint) error
}
