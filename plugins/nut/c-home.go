package nut

import (
	"net/http"
	"path"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"github.com/kapmahc/h2o/web"
)

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

	// favicon
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

// ------------

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
	err := p.Sitemap.ToXMLGz(p.Layout.Backend(c), c.Writer)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}

func (p *Plugin) getRssAtom(c *gin.Context) {
	lang := c.Param("lang")
	host := p.Layout.Backend(c)
	var author map[string]string
	if err := p.Settings.Get(p.DB, "site.author", &author); err != nil {
		author = map[string]string{
			"name":  "",
			"email": "",
		}
	}
	err := p.RSS.ToAtom(
		host,
		lang,
		p.I18n.T(lang, "site.title"),
		p.I18n.T(lang, "site.description"),
		&feeds.Author{
			Name:  author["name"],
			Email: author["email"],
		},
		c.Writer,
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
