package mail

import (
	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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
	return nil
}

func init() {
	web.Register(&Plugin{})
}
