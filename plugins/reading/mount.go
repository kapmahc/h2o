package reading

// Mount register
func (p *Plugin) Mount() error {

	ht := p.Router.Group("/reading/htdocs")
	ht.GET("/books", p.Layout.HTML("reading/books/index", p.indexBooksH))
	ht.GET("/books/:id", p.Layout.HTML("reading/books/show", p.showBookH))
	ht.GET("/pages/:id/*href", p.showPage)

	rt := p.Router.Group("/reading")
	rt.GET("/books", p.Layout.JSON(p.indexBooks))
	rt.DELETE("/books/:id", p.Layout.MustAdminMiddleware, p.Layout.JSON(p.destroyBook))
	return nil
}
