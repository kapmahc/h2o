package nut

import (
	"github.com/gin-gonic/gin"
	"github.com/kapmahc/h2o/web"
)

func (p *Plugin) getLocales(_ string, c *gin.Context) (interface{}, error) {
	items, err := p.I18n.All(c.Param("lang"))
	return items, err
}

func (p *Plugin) getSiteInfo(l string, c *gin.Context) (interface{}, error) {

	site := gin.H{}
	for _, k := range []string{"title", "subhead", "keywords", "description"} {
		site[k] = p.I18n.T(l, "site.title")
	}

	author := gin.H{}
	for _, k := range []string{"name", "email"} {
		var v string
		p.Settings.Get(p.DB, "site.author."+k, &v)
		author[k] = v
	}
	site["author"] = author

	langs, err := p.I18n.Languages(p.DB)
	if err != nil {
		return nil, err
	}

	site[web.LOCALE] = l
	site["languages"] = langs

	return site, nil
}
