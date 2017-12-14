package nut

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/kapmahc/h2o/web"
	log "github.com/sirupsen/logrus"
	"github.com/unrolled/render"
	validator "gopkg.in/go-playground/validator.v8"
)

const (
	// NOTICE notice
	NOTICE = "notice"
	// WARNING warning
	WARNING = "warning"
	// ERROR error
	ERROR = "error"

	// TITLE title
	TITLE = "title"
	// MESSAGE message
	MESSAGE = "message"

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
	Render *render.Render `inject:""`
	Dao    *Dao           `inject:""`
	Store  sessions.Store `inject:""`
	Jwt    *web.Jwt       `inject:""`
	DB     *gorm.DB       `inject:""`
	I18n   *web.I18n      `inject:""`
}

// MustSignInMiddleware must sign in middleware
func (p *Layout) MustSignInMiddleware(c *gin.Context) {

}

// MustAdminMiddleware must admin middleware
func (p *Layout) MustAdminMiddleware(c *gin.Context) {
	l := c.MustGet(web.LOCALE).(string)
	if _, ok := c.Get(CurrentUser); !ok {
		c.String(http.StatusUnauthorized, p.I18n.T(l, "errors.not-allow"))
		c.Abort()
		return
	}
}

// CurrentUserMiddleware current user middleware
func (p *Layout) CurrentUserMiddleware(c *gin.Context) {
	cm, err := p.Jwt.Parse(c.Request)
	if err != nil {
		log.Debug(err)
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

// Session get session
func (p *Layout) Session(c *gin.Context) *sessions.Session {
	ss, _ := p.Store.Get(c.Request, "session")
	return ss
}

// Save save session
func (p *Layout) Save(c *gin.Context, s *sessions.Session) {
	if err := s.Save(c.Request, c.Writer); err != nil {
		log.Error(err)
	}
}

// Redirect redirect
func (p *Layout) Redirect(to string, fn RedirectHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := fn(c.MustGet(web.LOCALE).(string), c); err != nil {
			log.Error(err)
			s := p.Session(c)
			s.AddFlash(err.Error(), ERROR)
			p.Save(c, s)
		}
		c.Redirect(http.StatusFound, to)
	}
}

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

// HTML wrap html
func (p *Layout) HTML(name string, handler HTMLHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.MustGet(web.LOCALE).(string)
		flashes := gin.H{}
		ss := p.Session(c)
		for _, n := range []string{NOTICE, WARNING, ERROR} {
			flashes[n] = ss.Flashes(n)
		}
		p.Save(c, ss)

		payload, err := handler(lang, c)
		if err != nil {
			payload = gin.H{}
		}
		payload["locale"] = lang
		payload["flashes"] = flashes
		payload["session"] = p.Session(c).Values
		if err == nil {
			p.Render.HTML(c.Writer, http.StatusOK, name, payload)
			return
		}
		log.Error(err)
		status, body := p.detectError(err)
		payload["reason"] = body
		payload["createdAt"] = time.Now()
		p.Render.HTML(c.Writer, status, "layouts/application/error", payload)
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
