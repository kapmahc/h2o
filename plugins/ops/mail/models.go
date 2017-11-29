package mail

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"time"
)

// https://www.linode.com/docs/email/postfix/email-with-postfix-dovecot-and-mysql
// http://wiki.dovecot.org/Authentication/PasswordSchemes
// https://mad9scientist.com/dovecot-password-creation-php/

// Domain domain
type Domain struct {
	ID        uint `gorm:"primary_key" json:"id"`
	Name      string
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// TableName table name
func (Domain) TableName() string {
	return "mail_domains"
}

// User user
type User struct {
	ID       uint `gorm:"primary_key" json:"id"`
	FullName string
	Email    string
	Password string
	Enable   bool

	UpdatedAt time.Time
	CreatedAt time.Time

	Domain   Domain
	DomainID uint
}

// TableName table name
func (User) TableName() string {
	return "mail_users"
}

func (p *User) sum(password string, salt []byte) string {
	buf := sha512.Sum512(append([]byte(password), salt...))
	return base64.StdEncoding.EncodeToString(append(buf[:], salt...))
}

// SetPassword set  password (SSHA512-CRYPT)
func (p *User) SetPassword(password string) error {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return err
	}
	p.Password = p.sum(password, salt)
	return nil
}

// ChkPassword check password
func (p *User) ChkPassword(password string) bool {
	buf, err := base64.StdEncoding.DecodeString(p.Password)
	if err != nil {
		return false
	}

	return len(buf) > sha512.Size && p.Password == p.sum(password, buf[sha512.Size:])
}

// Alias alias
type Alias struct {
	ID          uint `gorm:"primary_key" json:"id"`
	Source      string
	Destination string
	UpdatedAt   time.Time
	CreatedAt   time.Time

	Domain   Domain
	DomainID uint
}

// TableName table name
func (Alias) TableName() string {
	return "mail_aliases"
}
