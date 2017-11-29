package mall

import (
	"time"

	"github.com/kapmahc/h2o/plugins/nut"
)

// Address address
// http://www.bitboost.com/ref/international-address-formats.html#Formats
type Address struct {
	ID        uint
	Name      string
	Phone     string
	Zip       string
	Line1     string
	Line2     string
	City      string
	State     string
	Country   string
	CreatedAt time.Time
	UpdatedAt time.Time
	User      *nut.User
}

// Store store
type Store struct {
	ID          uint
	Name        string
	Description string
	Type        string
	Currency    string
	Metric      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Owner       *nut.User
	Address     *Address
}

// Tag tag
type Tag struct {
	ID          uint
	Name        string
	Description string
	Type        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Products    []*Product
}

// Vendor vendor
type Vendor struct {
	ID          uint
	Name        string
	Description string
	Type        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Store       *Store
}

// Product product
type Product struct {
	ID          uint
	Name        string
	Description string
	Type        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Store       *Store
	Vendor      *Vendor
	Tags        []*Tag
	Variants    []*Variant
}

// Variant variant
type Variant struct {
	ID          uint
	Sku         string
	Description string
	Type        string
	Price       float64
	Cost        float64
	Weight      float64
	Height      float64
	Width       float64
	Length      float64
	Stock       int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Product     *Product
}

// Journal journal
type Journal struct {
	ID        uint
	Action    string
	Quantity  int
	Variant   *Variant
	User      *nut.User
	CreatedAt time.Time
}

// Property property
type Property struct {
	ID        uint
	Key       string
	Val       string
	Variant   *Variant
	CreatedAt time.Time
	UpdatedAt time.Time
}
