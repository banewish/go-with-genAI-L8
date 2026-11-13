package model

type Category struct {
	ID        uint
	Name      string
	Books     []Book `gorm:"foreignKey:CategoryID"`
	CreatedAt string
	UpdatedAt string
}
