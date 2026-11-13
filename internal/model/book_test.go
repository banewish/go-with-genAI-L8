package model_test

import (
	"testing"
	"homework/internal/model"
)

func TestBook(t *testing.T) {
	book := &model.Book{
		ID:    1,
		Title: "Test Book",
	}
	if book.Title != "Test Book" {
		t.Errorf("Expected 'Test Book', got '%s'", book.Title)
	}
}

func TestCategory(t *testing.T) {
	category := &model.Category{
		ID:   1,
		Name: "Fiction",
	}
	if category.Name != "Fiction" {
		t.Errorf("Expected 'Fiction', got '%s'", category.Name)
	}
}

func TestAuthor(t *testing.T) {
	author := &model.Author{
		ID:       1,
		FullName: "John Doe",
	}
	if author.FullName != "John Doe" {
		t.Errorf("Expected 'John Doe', got '%s'", author.FullName)
	}
}
