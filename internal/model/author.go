package model

import "time"

type Author struct {
	ID        uint      `json:"id" gorm:"primaryKey;column:id"`
	FullName  string    `json:"full_name" gorm:"column:full_name"`
	Bio       *string   `json:"bio"`
	Books     []Book    `json:"books,omitempty" gorm:"many2many:book_authors"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}
