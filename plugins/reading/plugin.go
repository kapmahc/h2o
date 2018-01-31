package reading

import (
	"path"

	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/kapmahc/h2o/plugins/nut"
	"github.com/kapmahc/h2o/web"
)

// Plugin plugin
type Plugin struct {
	I18n    *web.I18n    `inject:""`
	Cache   *web.Cache   `inject:""`
	Sitemap *web.Sitemap `inject:""`
	Layout  *nut.Layout  `inject:""`
	Jwt     *web.Jwt     `inject:""`
	Router  *gin.Engine  `inject:""`
	DB      *gorm.DB     `inject:""`
}

// Init init beans
func (p *Plugin) Init(*inject.Graph) error {
	return nil
}

func (p *Plugin) root() string {
	return path.Join("tmp", "books")
}

func init() {
	web.Register(&Plugin{})
}
