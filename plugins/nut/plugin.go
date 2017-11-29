package nut

import (
	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/kapmahc/h2o/web"
	"github.com/unrolled/render"
)

// Plugin plugin
type Plugin struct {
	I18n     *web.I18n      `inject:""`
	Cache    *web.Cache     `inject:""`
	Jwt      *web.Jwt       `inject:""`
	Jobber   *web.Jobber    `inject:""`
	Settings *web.Settings  `inject:""`
	Router   *gin.Engine    `inject:""`
	DB       *gorm.DB       `inject:""`
	Render   *render.Render `inject:""`
}

// Init init beans
func (p *Plugin) Init(*inject.Graph) error {
	return nil
}

func init() {
	web.Register(&Plugin{})
}
