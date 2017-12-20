package forum

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kapmahc/h2o/plugins/nut"
	"github.com/kapmahc/h2o/web"
)

func (p *Plugin) indexComments(l string, c *gin.Context) (interface{}, error) {
	var items []Comment
	db := p.DB
	user := c.MustGet(nut.CurrentUser).(*nut.User)
	if c.MustGet(nut.IsAdmin).(bool) {
		db = p.DB.Where("lang = ?", l)
	} else {
		p.DB.Where("lang = ? AND user_id = ?", l, user.ID)
	}
	if err := db.Order("updated_at DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

type fmComment struct {
	Body string `json:"body" binding:"required"`
	Type string `json:"type" binding:"required"`
}

func (p *Plugin) createComment(l string, c *gin.Context) (interface{}, error) {
	user := c.MustGet(nut.CurrentUser).(*nut.User)
	var fm fmComment
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	it := Comment{
		Type:   fm.Type,
		Body:   fm.Body,
		UserID: user.ID,
	}
	if err := p.DB.Create(&it).Error; err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) showComment(l string, c *gin.Context) (interface{}, error) {
	var it = Comment{}
	if err := p.DB.Where("id = ?", c.Param("id")).First(&it).Error; err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) canEditComment(c *gin.Context) {
	user := c.MustGet(nut.CurrentUser).(*nut.User)
	var it = Comment{}
	if err := p.DB.Where("id = ?", c.Param("id")).First(&it).Error; err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}
	lang := c.MustGet(web.LOCALE).(string)
	if it.UserID != user.ID && !c.MustGet(nut.IsAdmin).(bool) {
		c.String(http.StatusForbidden, p.I18n.T(lang, "errors.not-allow"))
		c.Abort()
		return
	}
	c.Set("comment", &it)
}

func (p *Plugin) updateComment(l string, c *gin.Context) (interface{}, error) {
	var fm fmComment
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	it := c.MustGet("comment").(*Comment)

	if err := p.DB.Model(&it).Updates(map[string]interface{}{
		"body": fm.Body,
		"type": fm.Type,
	}).Error; err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

func (p *Plugin) destroyComment(l string, c *gin.Context) (interface{}, error) {
	it := c.MustGet("comment").(*Comment)
	if err := p.DB.Delete(it).Error; err != nil {
		return nil, err
	}
	return gin.H{}, nil
}
