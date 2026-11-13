package model

type Author struct {
	ID       uint
	FullName string
	Bio      string
	Books    []Book `gorm:"many2many:book_authors"`
	CreatedAt string
	UpdatedAt string
}
