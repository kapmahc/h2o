package nut

// Mount register
func (p *Plugin) Mount() error {
	p.Router.GET("/locales/:lang", p.Layout.JSON(p.getLocales))
	p.Router.GET("/site/info", p.Layout.JSON(p.getSiteInfo))
	return nil
}
