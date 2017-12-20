package survey

import (
	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/kapmahc/h2o/plugins/nut"
	"github.com/kapmahc/h2o/web"
	"github.com/unrolled/render"
	"github.com/urfave/cli"
)

// Plugin plugin
type Plugin struct {
	I18n   *web.I18n      `inject:""`
	Cache  *web.Cache     `inject:""`
	Jwt    *web.Jwt       `inject:""`
	Router *gin.Engine    `inject:""`
	DB     *gorm.DB       `inject:""`
	Render *render.Render `inject:""`
	Layout *nut.Layout    `inject:""`
}

// Init init beans
func (p *Plugin) Init(*inject.Graph) error {
	return nil
}

// Shell console commands
func (p *Plugin) Shell() []cli.Command {
	return []cli.Command{}
}

// Mount register
func (p *Plugin) Mount() error {
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

	return nil
}

func init() {
	web.Register(&Plugin{})
}
