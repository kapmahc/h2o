package nut

import (
	"fmt"
	"path"

	"github.com/ikeikeikeike/go-sitemap-generator/stm"
	"github.com/kapmahc/h2o/web"
	"github.com/spf13/viper"
)

func (p *Plugin) sitemap() ([]stm.URL, error) {
	var items []stm.URL
	for _, l := range p.Languages {
		items = append(items, stm.URL{
			"loc": fmt.Sprintf("/?locale=%s", l),
		})
	}
	return items, nil
}

// Mount register
func (p *Plugin) Mount() error {
	p.Sitemap.Register(p.sitemap)
	// --------------
	i18m, err := p.I18n.Middleware()
	if err != nil {
		return err
	}
	p.Router.Use(i18m)
	p.Router.Use(p.Layout.CurrentUserMiddleware)

	if web.MODE() != web.PRODUCTION {
		for k, v := range map[string]string{
			"3rd":    "node_modules",
			"assets": path.Join("themes", viper.GetString("server.theme"), "assets"),
		} {
			p.Router.Static("/"+k+"/", v)
		}
	}

	p.Router.GET("/", p.getHome)
	p.Router.GET("/robots.txt", p.getRobotsTxt)
	p.Router.GET("/sitemap.xml.gz", p.getSitemapGz)
	p.Router.GET("/locales/:lang", p.Layout.JSON(p.getLocales))
	p.Router.GET("/layout", p.Layout.JSON(p.getLayout))
	p.Router.POST("/leave-words", p.Layout.JSON(p.createLeaveWord))

	ung := p.Router.Group("/users")
	ung.POST("/sign-in", p.Layout.JSON(p.postUsersSignIn))
	ung.POST("/sign-up", p.Layout.JSON(p.postUsersSignUp))
	ung.POST("/confirm", p.Layout.JSON(p.postUsersConfirm))
	ung.POST("/unlock", p.Layout.JSON(p.postUsersUnlock))
	ung.POST("/forgot-password", p.Layout.JSON(p.postUsersForgotPassword))
	ung.POST("/reset-password", p.Layout.JSON(p.postUsersResetPassword))
	ung.GET("/confirm/:token", p.Layout.Redirect("/", p.getUsersConfirmToken))
	ung.GET("/unlock/:token", p.Layout.Redirect("/", p.getUsersUnlockToken))

	umg := p.Router.Group("/users", p.Layout.MustSignInMiddleware)
	umg.GET("/logs", p.Layout.JSON(p.getUsersLogs))
	umg.GET("/profile", p.Layout.JSON(p.getUsersProfile))
	umg.POST("/profile", p.Layout.JSON(p.postUsersProfile))
	umg.POST("/change-password", p.Layout.JSON(p.postUsersChangePassword))
	umg.DELETE("/sign-out", p.Layout.JSON(p.deleteUsersSignOut))

	p.Router.GET("/attachments", p.Layout.MustSignInMiddleware, p.Layout.JSON(p.indexAttachments))
	p.Router.POST("/attachments", p.Layout.MustSignInMiddleware, p.Layout.JSON(p.createAttachments))
	p.Router.DELETE("/attachments/:id", p.Layout.MustSignInMiddleware, p.Layout.JSON(p.destroyAttachments))

	ag := p.Router.Group("/admin", p.Layout.MustAdminMiddleware)
	ag.GET("/site/status", p.Layout.JSON(p.getAdminSiteStatus))
	ag.POST("/site/info", p.Layout.JSON(p.postAdminSiteInfo))
	ag.POST("/site/author", p.Layout.JSON(p.postAdminSiteAuthor))
	ag.GET("/site/seo", p.Layout.JSON(p.getAdminSiteSeo))
	ag.POST("/site/seo", p.Layout.JSON(p.postAdminSiteSeo))
	ag.GET("/site/smtp", p.Layout.JSON(p.getAdminSiteSMTP))
	ag.POST("/site/smtp", p.Layout.JSON(p.postAdminSiteSMTP))
	ag.PATCH("/site/smtp", p.Layout.JSON(p.patchAdminSiteSMTP))
	ag.POST("/site/home", p.Layout.JSON(p.postAdminSiteHome))
	ag.GET("/links", p.Layout.JSON(p.indexAdminLinks))
	ag.POST("/links", p.Layout.JSON(p.createAdminLink))
	ag.GET("/links/:id", p.Layout.JSON(p.showAdminLink))
	ag.POST("/links/:id", p.Layout.JSON(p.updateAdminLink))
	ag.DELETE("/links/:id", p.Layout.JSON(p.destroyAdminLink))
	ag.GET("/cards", p.Layout.JSON(p.indexAdminCards))
	ag.POST("/cards", p.Layout.JSON(p.createAdminCard))
	ag.GET("/cards/:id", p.Layout.JSON(p.showAdminCard))
	ag.POST("/cards/:id", p.Layout.JSON(p.updateAdminCard))
	ag.DELETE("/cards/:id", p.Layout.JSON(p.destroyAdminCard))
	ag.GET("/locales", p.Layout.JSON(p.indexAdminLocales))
	ag.POST("/locales", p.Layout.JSON(p.createAdminLocale))
	ag.GET("/locales/:id", p.Layout.JSON(p.showAdminLocale))
	ag.DELETE("/locales/:id", p.Layout.JSON(p.destroyAdminLocale))
	ag.GET("/friend-links", p.Layout.JSON(p.indexAdminFriendLinks))
	ag.POST("/friend-links", p.Layout.JSON(p.createAdminFriendLink))
	ag.GET("/friend-links/:id", p.Layout.JSON(p.showAdminFriendLink))
	ag.POST("/friend-links/:id", p.Layout.JSON(p.updateAdminFriendLink))
	ag.DELETE("/friend-links/:id", p.Layout.JSON(p.destroyAdminFriendLink))
	ag.GET("/leave-words", p.Layout.JSON(p.indexAdminLeaveWords))
	ag.DELETE("/leave-words/:id", p.Layout.JSON(p.destroyAdminLeaveWord))
	ag.GET("/users", p.Layout.JSON(p.indexAdminUsers))

	p.Jobber.Register(SendEmailJob, p.doSendEmail)
	return nil
}
