package nut

import (
	"path"

	"github.com/kapmahc/h2o/web"
	"github.com/spf13/viper"
)

// Mount register
func (p *Plugin) Mount() error {
	i18m, err := p.I18n.Middleware(p.DB)
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
	p.Router.GET("/locales/:lang", p.Layout.JSON(p.getLocales))
	p.Router.GET("/layout", p.Layout.JSON(p.getLayout))

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
	ag.POST("/site/home", p.Layout.JSON(p.postAdminSiteHome))

	p.Jobber.Register(SendEmailJob, p.doSendEmail)
	return nil
}
