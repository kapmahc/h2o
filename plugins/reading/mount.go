package reading

import (
	"fmt"

	"github.com/ikeikeikeike/go-sitemap-generator/stm"
)

// Mount register
func (p *Plugin) Mount() error {
	p.Sitemap.Register(p.sitemap)
	// ------------

	rt := p.Router.Group("/reading")
	rt.GET("/books", p.Layout.JSON(p.indexBooks))
	rt.GET("/books/:id", p.Layout.JSON(p.showBookH))
	rt.DELETE("/books/:id", p.Layout.MustAdminMiddleware, p.Layout.JSON(p.destroyBook))
	rt.GET("/notes", p.Layout.MustSignInMiddleware, p.Layout.JSON(p.indexNotes))
	rt.POST("/notes", p.Layout.MustSignInMiddleware, p.Layout.JSON(p.createNote))
	rt.GET("/notes/:id", p.Layout.JSON(p.showNote))
	rt.POST("/notes/:id", p.Layout.MustSignInMiddleware, p.canEditNote, p.Layout.JSON(p.updateNote))
	rt.DELETE("/notes/:id", p.Layout.MustSignInMiddleware, p.canEditNote, p.Layout.JSON(p.destroyNote))
	rt.GET("/pages/:id/*href", p.showPage)
	return nil
}

func (p *Plugin) sitemap() ([]stm.URL, error) {
	items := []stm.URL{
		{"loc": "/reading/books"},
		{"loc": "/reading/notes"},
	}

	var books []Book
	if err := p.DB.Select([]string{"id", "updated_at"}).Order("updated_at DESC").Find(&books).Error; err != nil {
		return nil, err
	}
	for _, it := range books {
		items = append(items, stm.URL{
			"loc":     fmt.Sprintf("/reading/books/%d", it.ID),
			"lastmod": it.UpdatedAt,
		})
	}

	return items, nil
}
