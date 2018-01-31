package survey

import (
	"fmt"

	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
	"github.com/ikeikeikeike/go-sitemap-generator/stm"
	"github.com/jinzhu/gorm"
	"github.com/kapmahc/h2o/plugins/nut"
	"github.com/kapmahc/h2o/web"
	"github.com/urfave/cli"
)

// Plugin plugin
type Plugin struct {
	I18n    *web.I18n    `inject:""`
	Cache   *web.Cache   `inject:""`
	Sitemap *web.Sitemap `inject:""`
	Jwt     *web.Jwt     `inject:""`
	Router  *gin.Engine  `inject:""`
	DB      *gorm.DB     `inject:""`
	Layout  *nut.Layout  `inject:""`
}

// Init init beans
func (p *Plugin) Init(*inject.Graph) error {
	return nil
}

// Shell console commands
func (p *Plugin) Shell() []cli.Command {
	return []cli.Command{}
}

func (p *Plugin) sitemap() ([]stm.URL, error) {
	items := []stm.URL{
		{"loc": "/survey/forms"},
	}

	var forms []Form
	if err := p.DB.Select([]string{"id", "updated_at"}).Order("updated_at DESC").Find(&forms).Error; err != nil {
		return nil, err
	}
	for _, it := range forms {
		items = append(items, stm.URL{
			"loc":     fmt.Sprintf("/survey/forms/apply/%d", it.ID),
			"lastmod": it.UpdatedAt,
		})
	}

	return items, nil
}

// Mount register
func (p *Plugin) Mount() error {
	p.Sitemap.Register(p.sitemap)

	rt := p.Router.Group("/survey")
	rt.GET("/forms", p.Layout.MustSignInMiddleware, p.Layout.JSON(p.indexForms))
	rt.POST("/forms", p.Layout.MustSignInMiddleware, p.Layout.JSON(p.createForm))
	rt.GET("/forms/:id", p.Layout.JSON(p.showForm))
	rt.POST("/forms/:id", p.Layout.MustSignInMiddleware, p.canEditForm, p.Layout.JSON(p.updateForm))
	rt.DELETE("/forms/:id", p.Layout.MustSignInMiddleware, p.canEditForm, p.Layout.JSON(p.destroyForm))

	rt.GET("/fields", p.Layout.JSON(p.indexFields))
	rt.POST("/fields", p.Layout.MustSignInMiddleware, p.Layout.JSON(p.createField))
	rt.GET("/fields/:id", p.Layout.JSON(p.showField))
	rt.POST("/fields/:id", p.Layout.MustSignInMiddleware, p.canEditField, p.Layout.JSON(p.updateField))
	rt.DELETE("/fields/:id", p.Layout.MustSignInMiddleware, p.canEditField, p.Layout.JSON(p.destroyField))

	rt.GET("/records", p.Layout.MustSignInMiddleware, p.Layout.JSON(p.indexRecords))
	rt.DELETE("/records/:id", p.Layout.MustSignInMiddleware, p.canEditRecord, p.Layout.JSON(p.destroyRecord))

	return nil
}

func init() {
	web.Register(&Plugin{})
}
