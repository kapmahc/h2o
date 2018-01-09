package reading

import (
	"fmt"

	"github.com/ikeikeikeike/go-sitemap-generator/stm"
)

// Mount register
func (p *Plugin) Mount() error {
	p.Sitemap.Register(p.sitemap)
	// ------------

	ht := p.Router.Group("/reading/htdocs")
	ht.GET("/notes", p.Layout.HTML("reading/notes/index", p.indexNotesH))
	ht.GET("/books", p.Layout.HTML("reading/books/index", p.indexBooksH))
	ht.GET("/books/:id", p.Layout.HTML("reading/books/show", p.showBookH))
	ht.GET("/pages/:id/*href", p.showPage)

	rt := p.Router.Group("/reading")
	rt.GET("/books", p.Layout.JSON(p.indexBooks))
	rt.DELETE("/books/:id", p.Layout.MustAdminMiddleware, p.Layout.JSON(p.destroyBook))
	rt.GET("/notes", p.Layout.MustSignInMiddleware, p.Layout.JSON(p.indexNotes))
	rt.POST("/notes", p.Layout.MustSignInMiddleware, p.Layout.JSON(p.createNote))
	rt.GET("/notes/:id", p.Layout.JSON(p.showNote))
	rt.POST("/notes/:id", p.Layout.MustSignInMiddleware, p.canEditNote, p.Layout.JSON(p.updateNote))
	rt.DELETE("/notes/:id", p.Layout.MustSignInMiddleware, p.canEditNote, p.Layout.JSON(p.destroyNote))
	return nil
}

func (p *Plugin) sitemap() ([]stm.URL, error) {
	items := []stm.URL{
		{"loc": "/reading/htdocs/books"},
		{"loc": "/reading/htdocs/notes"},
	}

	var books []Book
	if err := p.DB.Select([]string{"id", "updated_at"}).Order("updated_at DESC").Find(&books).Error; err != nil {
		return nil, err
	}
	for _, it := range books {
		items = append(items, stm.URL{
			"loc":     fmt.Sprintf("/reading/htdocs/books/%d", it.ID),
			"lastmod": it.UpdatedAt,
		})
	}

	return items, nil
}
