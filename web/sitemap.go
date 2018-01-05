package web

import (
	"compress/gzip"
	"io"

	"github.com/ikeikeikeike/go-sitemap-generator/stm"
)

// SitemapHandler sitemap handler
type SitemapHandler func() ([]stm.URL, error)

// NewSitemap new sitemap
func NewSitemap() *Sitemap {
	return &Sitemap{
		handlers: make([]SitemapHandler, 0),
	}
}

// Sitemap sitemap helper
type Sitemap struct {
	handlers []SitemapHandler
}

// Register register handler
func (p *Sitemap) Register(handlers ...SitemapHandler) {
	p.handlers = append(p.handlers, handlers...)
}

// Generate xml.gz content body
func (p *Sitemap) Generate(h string, w io.Writer) error {
	sm := stm.NewSitemap()
	sm.Create()
	sm.SetDefaultHost(h)
	for _, hnd := range p.handlers {
		items, err := hnd()
		if err != nil {
			return err
		}
		for _, it := range items {
			sm.Add(it)
		}
	}
	buf := sm.XMLContent()

	wrt := gzip.NewWriter(w)
	defer wrt.Close()
	wrt.Write(buf)
	return nil
}
