package nut

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/kapmahc/h2o/web"
	log "github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v8"
)

const (
	// UID uid
	UID = "uid"
	// CurrentUser current user
	CurrentUser = "current-user"
	// IsAdmin is admin?
	IsAdmin = "is-admin"
)

// HTMLHandlerFunc html handler func
type HTMLHandlerFunc func(string, *gin.Context) (gin.H, error)

// RedirectHandlerFunc redirect handle func
type RedirectHandlerFunc func(string, *gin.Context) error

// ObjectHandlerFunc object handle func
type ObjectHandlerFunc func(string, *gin.Context) (interface{}, error)

// Layout layout
type Layout struct {
	Dao       *Dao      `inject:""`
	Jwt       *web.Jwt  `inject:""`
	DB        *gorm.DB  `inject:""`
	I18n      *web.I18n `inject:""`
	Languages []string  `inject:"languages"`
}

// MustSignInMiddleware must sign in middleware
func (p *Layout) MustSignInMiddleware(c *gin.Context) {
	l := c.MustGet(web.LOCALE).(string)
	if _, ok := c.Get(CurrentUser); ok {
		return
	}
	c.String(http.StatusForbidden, p.I18n.T(l, "errors.not-allow"))
	c.Abort()
}

// MustAdminMiddleware must admin middleware
func (p *Layout) MustAdminMiddleware(c *gin.Context) {
	l := c.MustGet(web.LOCALE).(string)
	if is, ok := c.Get(IsAdmin); ok && is.(bool) {
		return
	}
	c.String(http.StatusForbidden, p.I18n.T(l, "errors.not-allow"))
	c.Abort()
}

// CurrentUserMiddleware current user middleware
func (p *Layout) CurrentUserMiddleware(c *gin.Context) {
	cm, err := p.Jwt.Parse(c.Request)
	if err != nil {
		log.Error(err)
		return
	}
	uid, ok := cm.Get(UID).(string)
	if !ok {
		return
	}
	user, err := p.Dao.GetUserByUID(p.DB, uid)
	if err != nil {
		log.Error(err)
		return
	}
	if !user.IsConfirm() || user.IsLock() {
		return
	}
	c.Set(CurrentUser, user)
	c.Set(IsAdmin, p.Dao.Is(p.DB, user.ID, RoleAdmin))
}

// Redirect redirect
// func (p *Layout) Redirect(to string, fn RedirectHandlerFunc) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if err := fn(c.MustGet(web.LOCALE).(string), c); err != nil {
// 			log.Error(err)
// 			c.String(http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 		c.Redirect(http.StatusFound, to)
// 	}
// }

// JSON render json
func (p *Layout) JSON(fn ObjectHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if val, err := fn(c.MustGet(web.LOCALE).(string), c); err == nil {
			c.JSON(http.StatusOK, val)
		} else {
			log.Error(err)
			status, body := p.detectError(err)
			c.String(status, body)
		}
	}
}

func (p *Layout) detectError(e error) (int, string) {
	if er, ok := e.(validator.ValidationErrors); ok {
		var ss []string
		for _, it := range er {
			ss = append(ss, fmt.Sprintf("Validation for '%s' failed on the '%s' tag;", it.Field, it.Tag))
		}
		return http.StatusBadRequest, strings.Join(ss, "\n")
	}
	return http.StatusInternalServerError, e.Error()
}

// XML wrap xml
func (p *Layout) XML(fn ObjectHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if val, err := fn(c.MustGet(web.LOCALE).(string), c); err == nil {
			c.XML(http.StatusOK, val)
		} else {
			log.Error(err)
			status, body := p.detectError(err)
			c.String(status, body)
		}
	}
}

// Backend backend home url
func (p *Layout) Backend(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme += "s"
	}
	return scheme + "://" + c.Request.Host
}

// Frontend frontend home url
func (p *Layout) Frontend(c *gin.Context) string {
	return c.GetHeader("Origin")
}
