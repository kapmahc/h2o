package nut

import (
	"bytes"
	h_template "html/template"
	"net/http"
	"path"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"github.com/kapmahc/h2o/web"
)

func (p *Plugin) getHome(c *gin.Context) {
	lang := c.MustGet(web.LOCALE).(string)
	// carousel off-canvas
	theme := c.Query("theme")
	var home map[string]string
	if err := p.Settings.Get(p.DB, "site.home."+lang, &home); err == nil {
		href := home["href"]
		if href != "" {
			c.Redirect(http.StatusFound, href)
			return
		}
		if theme == "" {
			theme = home["theme"]
		}
	}
	if theme == "" {
		theme = "off-canvas"
	}

	p.Layout.HTML("nut/home/"+theme, func(string, *gin.Context) (gin.H, error) {
		return gin.H{}, nil
	})(c)
}

func (p *Plugin) getDonate(l string, c *gin.Context) (gin.H, error) {
	item := make(map[string]interface{})
	if err := p.Settings.Get(p.DB, "site.donate", &item); err != nil {
		return nil, err
	}
	rst := gin.H{
		"body": item["body"],
		"type": item["type"],
		TITLE:  p.I18n.T(l, "nut.donate.title"),
	}
	if paypal, ok := item["paypal"]; ok && paypal.(string) != "" {
		var buf bytes.Buffer
		tpl, err := template.ParseFiles(path.Join("templates", "paypal.html"))
		if err != nil {
			return nil, err
		}
		if err := tpl.Execute(&buf, gin.H{"id": paypal}); err != nil {
			return nil, err
		}
		rst["paypal"] = h_template.HTML(buf.Bytes())
	}
	return rst, nil
}

// ------------------

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
