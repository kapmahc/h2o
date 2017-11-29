package reading

import (
	"time"

	"github.com/kapmahc/h2o/plugins/nut"
)

// Book book
type Book struct {
	ID uint `gorm:"primary_key" json:"id"`

	Author      string    `json:"author"`
	Publisher   string    `json:"publisher"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	Lang        string    `json:"lang"`
	File        string    `json:"-"`
	Subject     string    `json:"subject"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"publishedAt"`
	Cover       string    `json:"cover"`
	UpdatedAt   time.Time
	CreatedAt   time.Time
}

// TableName table name
func (Book) TableName() string {
	return "reading_books"
}

// Note note
type Note struct {
	ID        uint `gorm:"primary_key" json:"id"`
	Type      string
	Body      string
	UpdatedAt time.Time
	CreatedAt time.Time

	User   nut.User
	UserID uint
	Book   Book
	BookID uint
}

// TableName table name
func (Note) TableName() string {
	return "reading_notes"
}
