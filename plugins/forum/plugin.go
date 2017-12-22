package forum

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
	Layout *nut.Layout    `inject:""`
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
	ht := p.Router.Group("/forum/htdocs")
	ht.GET("/tags", p.Layout.HTML("forum/tags/index", p.indexTagsH))
	ht.GET("/tags/:id", p.Layout.HTML("forum/tags/show", p.showTagH))
	ht.GET("/articles", p.Layout.HTML("forum/articles/index", p.indexArticlesH))
	ht.GET("/articles/:id", p.Layout.HTML("forum/articles/show", p.showArticleH))
	ht.GET("/comments", p.Layout.HTML("forum/comments/index", p.indexCommentsH))

	rt := p.Router.Group("/forum")

	rt.GET("/articles", p.Layout.MustSignInMiddleware, p.Layout.JSON(p.indexArticles))
	rt.POST("/articles", p.Layout.MustSignInMiddleware, p.Layout.JSON(p.createArticle))
	rt.GET("/articles/:id", p.Layout.JSON(p.showArticle))
	rt.POST("/articles/:id", p.Layout.MustSignInMiddleware, p.canEditArticle, p.Layout.JSON(p.updateArticle))
	rt.DELETE("/articles/:id", p.Layout.MustSignInMiddleware, p.canEditArticle, p.Layout.JSON(p.destroyArticle))

	rt.GET("/comments", p.Layout.MustSignInMiddleware, p.Layout.JSON(p.indexComments))
	rt.POST("/comments", p.Layout.MustSignInMiddleware, p.Layout.JSON(p.createComment))
	rt.GET("/comments/:id", p.Layout.JSON(p.showComment))
	rt.POST("/comments/:id", p.Layout.MustSignInMiddleware, p.canEditComment, p.Layout.JSON(p.updateComment))
	rt.DELETE("/comments/:id", p.Layout.MustSignInMiddleware, p.canEditComment, p.Layout.JSON(p.destroyComment))

	rt.GET("/tags", p.Layout.MustSignInMiddleware, p.Layout.JSON(p.indexTags))
	rt.POST("/tags", p.Layout.MustAdminMiddleware, p.Layout.JSON(p.createTag))
	rt.GET("/tags/:id", p.Layout.JSON(p.showTag))
	rt.POST("/tags/:id", p.Layout.MustSignInMiddleware, p.Layout.MustAdminMiddleware, p.Layout.JSON(p.updateTag))
	rt.DELETE("/tags/:id", p.Layout.MustSignInMiddleware, p.Layout.MustAdminMiddleware, p.Layout.JSON(p.destroyTag))
	return nil
}

func init() {
	web.Register(&Plugin{})
}
