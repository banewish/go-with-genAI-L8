package model

type Copy struct {
	ID     uint
	BookID uint
	UserID uint
	Book   Book
	User   User
}
