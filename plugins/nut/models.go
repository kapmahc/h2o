package nut

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	// RoleAdmin admin role
	RoleAdmin = "admin"
	// RoleRoot root role
	RoleRoot = "root"
	// UserTypeEmail email user
	UserTypeEmail = "email"

	// DefaultResourceType default resource type
	DefaultResourceType = ""
	// DefaultResourceID default resourc id
	DefaultResourceID = 0
)

// User user
type User struct {
	ID              uint       `json:"id" gorm:"primary_key"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	UID             string     `json:"uid" gorm:"column:uid"`
	Password        []byte     `json:"-"`
	ProviderID      string     `json:"providerId"`
	ProviderType    string     `json:"providerType"`
	Logo            string     `json:"logo"`
	SignInCount     uint       `json:"signInCount"`
	LastSignInAt    *time.Time `json:"lastSignInAt"`
	LastSignInIP    string     `json:"lastSignInIp"`
	CurrentSignInAt *time.Time `json:"currentSignInAt"`
	CurrentSignInIP string     `json:"currentSignInIp"`
	ConfirmedAt     *time.Time `json:"confirmedAt"`
	LockedAt        *time.Time `json:"lockAt"`
	Logs            []Log      `json:"logs"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	CreatedAt       time.Time  `json:"createdAt"`
}

// IsConfirm is confirm?
func (p *User) IsConfirm() bool {
	return p.ConfirmedAt != nil
}

// IsLock is lock?
func (p *User) IsLock() bool {
	return p.LockedAt != nil
}

// SetGravatarLogo set logo by gravatar
func (p *User) SetGravatarLogo() {
	// https: //en.gravatar.com/site/implement/
	buf := md5.Sum([]byte(strings.ToLower(p.Email)))
	p.Logo = fmt.Sprintf("https://gravatar.com/avatar/%s.png", hex.EncodeToString(buf[:]))
}

//SetUID generate uid
func (p *User) SetUID() {
	p.UID = uuid.New().String()
}

func (p User) String() string {
	return fmt.Sprintf("%s<%s>", p.Name, p.Email)
}

// TableName table name
func (p User) TableName() string {
	return "users"
}

// Attachment attachment
type Attachment struct {
	ID           uint      `json:"id" gorm:"primary_key"`
	Title        string    `json:"title"`
	URL          string    `json:"url"`
	Length       int64     `json:"length"`
	MediaType    string    `json:"mediaType"`
	ResourceID   uint      `json:"resourceId" sql:",notnull"`
	ResourceType string    `json:"resourceType" sql:",notnull"`
	User         User      `json:"user"`
	UserID       uint      `json:"userId"`
	UpdatedAt    time.Time `json:"updatedAt"`
	CreatedAt    time.Time `json:"crateAt"`
}

// IsPicture is picture?
func (p *Attachment) IsPicture() bool {
	return strings.HasPrefix(p.MediaType, "image/")
}

// TableName table name
func (p Attachment) TableName() string {
	return "attachments"
}

// Log log
type Log struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Message   string    `json:"message"`
	IP        string    `json:"ip"`
	User      User      `json:"user"`
	UserID    uint      `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

func (p Log) String() string {
	return fmt.Sprintf("%s: [%s]\t %s", p.CreatedAt.Format(time.ANSIC), p.IP, p.Message)
}

// TableName table name
func (p Log) TableName() string {
	return "logs"
}

// Policy policy
type Policy struct {
	ID        uint `gorm:"primary_key"`
	Begin     time.Time
	StartUp   time.Time
	ShutDown  time.Time
	User      User
	UserID    uint
	Role      Role
	RoleID    uint
	UpdatedAt time.Time
	CreatedAt time.Time
}

//Enable is enable?
func (p *Policy) Enable() bool {
	now := time.Now()
	return now.After(p.StartUp) && now.Before(p.ShutDown)
}

// TableName table name
func (p Policy) TableName() string {
	return "policies"
}

// Role role
type Role struct {
	ID           uint `gorm:"primary_key"`
	Name         string
	ResourceID   uint
	ResourceType string
	UpdatedAt    time.Time
	CreatedAt    time.Time
}

func (p Role) String() string {
	return fmt.Sprintf("%s@%s://%d", p.Name, p.ResourceType, p.ResourceID)
}

// TableName table name
func (p Role) TableName() string {
	return "roles"
}

// Vote vote
type Vote struct {
	ID           uint `gorm:"primary_key"`
	Point        int
	ResourceID   uint
	ResourceType string
	UpdatedAt    time.Time
	CreatedAt    time.Time
}

// TableName table name
func (p Vote) TableName() string {
	return "votes"
}

// LeaveWord leave-word
type LeaveWord struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
}

// TableName table name
func (p LeaveWord) TableName() string {
	return "leave_words"
}

// Link link
type Link struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Lang      string    `json:"lang"`
	Loc       string    `json:"loc"`
	Href      string    `json:"href"`
	Label     string    `json:"label"`
	SortOrder int       `json:"sortOrder"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// TableName table name
func (p Link) TableName() string {
	return "links"
}

// Card card
type Card struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Lang      string    `json:"lang"`
	Loc       string    `json:"loc"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	Type      string    `json:"type"`
	Href      string    `json:"href"`
	Logo      string    `json:"logo"`
	SortOrder int       `json:"sortOrder" sql:",notnull"`
	Action    string    `json:"action"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// TableName table name
func (p Card) TableName() string {
	return "cards"
}

// FriendLink friend_links
type FriendLink struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Title     string    `json:"title"`
	Home      string    `json:"home"`
	Logo      string    `json:"logo"`
	SortOrder int       `json:"sortOrder"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// TableName table name
func (p FriendLink) TableName() string {
	return "friend_links"
}
