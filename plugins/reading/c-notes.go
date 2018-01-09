package reading

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kapmahc/h2o/plugins/nut"
	"github.com/kapmahc/h2o/web"
)

func (p *Plugin) indexNotes(l string, c *gin.Context) (interface{}, error) {
	var items []Note
	db := p.DB
	user := c.MustGet(nut.CurrentUser).(*nut.User)
	if !c.MustGet(nut.IsAdmin).(bool) {
		p.DB.Where("user_id = ?", user.ID)
	}
	if err := db.Order("updated_at DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

type fmNote struct {
	Body string `json:"body" binding:"required"`
	Type string `json:"type" binding:"required"`
}

func (p *Plugin) createNote(l string, c *gin.Context) (interface{}, error) {
	user := c.MustGet(nut.CurrentUser).(*nut.User)
	var fm fmNote
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	bid, err := strconv.Atoi(c.Query("bookId"))
	if err != nil {
		return nil, err
	}
	it := Note{
		Type:   fm.Type,
		Body:   fm.Body,
		BookID: uint(bid),
		UserID: user.ID,
	}
	if err := p.DB.Create(&it).Error; err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) showNote(l string, c *gin.Context) (interface{}, error) {
	var it = Note{}
	if err := p.DB.Where("id = ?", c.Param("id")).First(&it).Error; err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) canEditNote(c *gin.Context) {
	user := c.MustGet(nut.CurrentUser).(*nut.User)
	var it = Note{}
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
	c.Set("note", &it)
}

func (p *Plugin) updateNote(l string, c *gin.Context) (interface{}, error) {
	var fm fmNote
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	it := c.MustGet("note").(*Note)

	if err := p.DB.Model(&it).Updates(map[string]interface{}{
		"body": fm.Body,
		"type": fm.Type,
	}).Error; err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

func (p *Plugin) destroyNote(l string, c *gin.Context) (interface{}, error) {
	it := c.MustGet("note").(*Note)
	if err := p.DB.Delete(it).Error; err != nil {
		return nil, err
	}
	return gin.H{}, nil
}
