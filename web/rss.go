package web

import (
	"fmt"
	"io"
	"time"

	"github.com/gorilla/feeds"
)

// NewRSS new RSS
func NewRSS() *RSS {
	return &RSS{
		handlers: make([]RSSHandler, 0),
	}
}

// RSSHandler rss handler
type RSSHandler func(l string) ([]*feeds.Item, error)

// RSS rss helper
type RSS struct {
	handlers []RSSHandler
}

// Register register handler
func (p *RSS) Register(handlers ...RSSHandler) {
	p.handlers = append(p.handlers, handlers...)
}

// ToAtom write to atom
func (p *RSS) ToAtom(host, lang, title, dest string, author *feeds.Author, wrt io.Writer) error {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       title,
		Link:        &feeds.Link{Href: fmt.Sprintf("%s/?locale=%s", host, lang)},
		Description: dest,
		Author:      author,
		Created:     now,
		Items:       make([]*feeds.Item, 0),
	}

	for _, hnd := range p.handlers {
		items, err := hnd(lang)
		if err != nil {
			return err
		}
		feed.Items = append(feed.Items, items...)
	}
	return feed.WriteAtom(wrt)
}
