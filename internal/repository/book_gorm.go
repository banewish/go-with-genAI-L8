package repository

import (
	"context"
	"errors"
	"time"

	"homework/internal/model"

	"gorm.io/gorm"
)

type GormBookRepository struct {
	db *gorm.DB
}

func NewGormBookRepository(db *gorm.DB) BookRepository {
	return &GormBookRepository{db: db}
}

func (r *GormBookRepository) Create(ctx context.Context, req *model.BookRequest) (*model.Book, error) {
	book := &model.Book{
		Title:       req.Title,
		ISBN:        req.ISBN,
		CategoryID:  req.CategoryID,
		Description: req.Description,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if req.PublishedDate != nil {
		if t, err := time.Parse("2006-01-02", *req.PublishedDate); err == nil {
			book.PublishedDate = &t
		}
	}

	if err := r.db.WithContext(ctx).Create(book).Error; err != nil {
		return nil, err
	}

	if len(req.AuthorIDs) > 0 {
		authors := make([]model.Author, len(req.AuthorIDs))
		for i, authorID := range req.AuthorIDs {
			authors[i].ID = authorID
		}
		if err := r.db.WithContext(ctx).Model(book).Association("Authors").Append(authors); err != nil {
			return nil, err
		}
	}

	return r.FindByID(ctx, book.ID)
}

func (r *GormBookRepository) FindByID(ctx context.Context, id uint) (*model.Book, error) {
	book := &model.Book{}

	if err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Authors").
		First(book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return book, nil
}

func (r *GormBookRepository) List(ctx context.Context, offset, limit int) ([]model.Book, error) {
	var books []model.Book

	// AFTER FIX: Fetch books with category JOIN in one query
	if err := r.db.WithContext(ctx).
		Joins("LEFT JOIN categories ON books.category_id = categories.id").
		Offset(offset).
		Limit(limit).
		Order("books.id ASC").
		Find(&books).Error; err != nil {
		return nil, err
	}

	return books, nil
}

func (r *GormBookRepository) Update(ctx context.Context, id uint, req *model.BookRequest) (*model.Book, error) {
	book, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if book == nil {
		return nil, gorm.ErrRecordNotFound
	}

	updates := map[string]interface{}{
		"title":       req.Title,
		"isbn":        req.ISBN,
		"category_id": req.CategoryID,
		"description": req.Description,
		"updated_at":  time.Now().UTC(),
	}

	if req.PublishedDate != nil {
		if t, err := time.Parse("2006-01-02", *req.PublishedDate); err == nil {
			updates["published_date"] = t
		}
	}

	if err := r.db.WithContext(ctx).Model(book).Updates(updates).Error; err != nil {
		return nil, err
	}

	if len(req.AuthorIDs) > 0 {
		if err := r.db.WithContext(ctx).Model(book).Association("Authors").Clear(); err != nil {
			return nil, err
		}

		authors := make([]model.Author, len(req.AuthorIDs))
		for i, authorID := range req.AuthorIDs {
			authors[i].ID = authorID
		}
		if err := r.db.WithContext(ctx).Model(book).Association("Authors").Append(authors); err != nil {
			return nil, err
		}
	}

	return r.FindByID(ctx, id)
}

func (r *GormBookRepository) Delete(ctx context.Context, id uint) error {
	book := &model.Book{}
	if err := r.db.WithContext(ctx).Delete(book, id).Error; err != nil {
		return err
	}
	return nil
}
