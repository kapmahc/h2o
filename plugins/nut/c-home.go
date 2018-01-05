package nut

import (
	"net/http"
	"path"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/kapmahc/h2o/web"
)

func (p *Plugin) getHome(c *gin.Context) {
	// carousel off-canvas
	theme := c.Query("theme")
	if theme == "" {
		if err := p.Settings.Get(p.DB, "site.home.theme", &theme); err != nil {
			theme = "off-canvas"
		}
	}
	p.Layout.HTML("nut/home/"+theme, func(string, *gin.Context) (gin.H, error) {
		return gin.H{}, nil
	})(c)
}

func (p *Plugin) getLocales(_ string, c *gin.Context) (interface{}, error) {
	items, err := p.I18n.All(c.Param("lang"))
	return items, err
}

func (p *Plugin) getLayout(l string, c *gin.Context) (interface{}, error) {
	// site info
	site := gin.H{}
	for _, k := range []string{"title", "subhead", "keywords", "description", "copyright"} {
		site[k] = p.I18n.T(l, "site."+k)
	}
	author := make(map[string]string)
	p.Settings.Get(p.DB, "site.author", &author)
	site["author"] = author

	// home
	var home string
	p.Settings.Get(p.DB, "site.home.theme", &home)
	site["home"] = home
	var favicon string
	p.Settings.Get(p.DB, "site.favicon", &favicon)
	site["favicon"] = favicon

	// i18n

	site[web.LOCALE] = l
	site["languages"] = p.Languages[:]

	// current-user
	user, ok := c.Get(CurrentUser)
	// nav
	if ok {
		user := user.(*User)
		site["user"] = gin.H{
			"name":  user.Name,
			"type":  user.ProviderType,
			"admin": c.MustGet(IsAdmin).(bool),
		}
	}

	return site, nil
}

// http://www.robotstxt.org/robotstxt.html
func (p *Plugin) getRobotsTxt(c *gin.Context) {
	tpl, err := template.ParseFiles(path.Join("templates", "robots.txt"))
	if err == nil {
		if err = tpl.Execute(c.Writer, gin.H{"home": p.Layout.Backend(c)}); err == nil {
			return
		}
	}
	c.String(http.StatusInternalServerError, err.Error())
}

func (p *Plugin) getSitemapGz(c *gin.Context) {
	err := p.Sitemap.Generate(p.Layout.Backend(c), c.Writer)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
