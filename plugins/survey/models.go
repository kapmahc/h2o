package survey

import (
	"time"

	"github.com/kapmahc/h2o/plugins/nut"
)

// Form form
type Form struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Mode      string    `json:"mode"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	StartUp   time.Time `json:"startUp"`
	ShutDown  time.Time `json:"shutDown"`
	UpdatedAt time.Time
	CreatedAt time.Time

	User    nut.User
	UserID  uint
	Fields  []Field  `json:"fields"`
	Records []Record `json:"records"`
}

// TableName table name
func (Form) TableName() string {
	return "survey_forms"
}

// Available available?
func (p *Form) Available() bool {
	now := time.Now()
	return now.After(p.StartUp) && now.Before(p.ShutDown)
}

// Field field
type Field struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	Label     string    `json:"label"`
	Value     string    `json:"value"`
	Required  bool      `json:"required"`
	SortOrder int       `json:"sortOrder"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`

	Form   Form
	FormID uint
}

// TableName table name
func (Field) TableName() string {
	return "survey_fields"
}

// Record record
type Record struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Value     string    `json:"value"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`

	Form   Form
	FormID uint
}

// TableName table name
func (Record) TableName() string {
	return "survey_records"
}
