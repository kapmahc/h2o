package nut

// Mount register
func (p *Plugin) Mount() error {
	p.Router.GET("/", p.getHome)
	p.Router.GET("/locales/:lang", p.Layout.JSON(p.getLocales))
	p.Router.GET("/layout", p.Layout.JSON(p.getLayout))

	p.Router.POST("/users/sign-in", p.Layout.JSON(p.postUsersSignIn))
	p.Router.POST("/users/sign-up", p.Layout.JSON(p.postUsersSignUp))
	p.Router.POST("/users/confirm", p.Layout.JSON(p.postUsersConfirm))
	p.Router.POST("/users/unlock", p.Layout.JSON(p.postUsersUnlock))
	p.Router.POST("/users/forgot-password", p.Layout.JSON(p.postUsersForgotPassword))
	p.Router.GET("/users/confirm/:token", p.Layout.Redirect("/", p.getUsersConfirmToken))
	p.Router.GET("/users/unlock/:token", p.Layout.Redirect("/", p.getUsersUnlockToken))

	p.Jobber.Register(SendEmailJob, p.doSendEmail)
	return nil
}
