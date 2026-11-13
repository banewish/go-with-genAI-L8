package model

import "time"

type Book struct {
	ID            uint       `json:"id" gorm:"primaryKey;column:id"`
	Title         string     `json:"title"`
	ISBN          *string    `json:"isbn"`
	CategoryID    uint       `json:"category_id" gorm:"column:category_id;index"`
	Category      Category   `json:"category" gorm:"foreignKey:CategoryID"`
	PublishedDate *time.Time `json:"published_date" gorm:"column:published_date"`
	Description   *string    `json:"description"`
	Authors       []Author   `json:"authors,omitempty" gorm:"many2many:book_authors"`
	CreatedAt     time.Time  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"column:updated_at"`
}

type BookRequest struct {
	Title         string  `json:"title"`
	ISBN          *string `json:"isbn"`
	CategoryID    uint    `json:"category_id"`
	PublishedDate *string `json:"published_date"`
	Description   *string `json:"description"`
	AuthorIDs     []uint  `json:"author_ids"`
}

type AuthorRequest struct {
	FullName string  `json:"full_name"`
	Bio      *string `json:"bio"`
}

type CategoryRequest struct {
	Name string `json:"name"`
}
