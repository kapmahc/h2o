package survey

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kapmahc/h2o/plugins/nut"
	"github.com/kapmahc/h2o/web"
)

func (p *Plugin) indexForms(l string, c *gin.Context) (interface{}, error) {
	var items []Form
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

type fmForm struct {
	Title    string `json:"title" binding:"required"`
	Body     string `json:"body" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Mode     string `json:"mode" binding:"required"`
	StartUp  string `json:"startUp" binding:"required"`
	ShutDown string `json:"shutDown" binding:"required"`
}

func (p *fmForm) parse() (*time.Time, *time.Time, error) {
	begin, err := time.Parse(web.DateFormat, p.StartUp)
	if err != nil {
		return nil, nil, err
	}
	end, err := time.Parse(web.DateFormat, p.ShutDown)
	if err != nil {
		return nil, nil, err
	}
	return &begin, &end, nil
}

func (p *Plugin) createForm(l string, c *gin.Context) (interface{}, error) {
	user := c.MustGet(nut.CurrentUser).(*nut.User)
	var fm fmForm
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}

	begin, end, err := fm.parse()
	if err != nil {
		return nil, err
	}

	it := Form{
		Title:    fm.Title,
		Body:     fm.Body,
		Type:     fm.Type,
		Mode:     fm.Mode,
		StartUp:  *begin,
		ShutDown: *end,
		UserID:   user.ID,
	}

	if err := p.DB.Create(&it).Error; err != nil {
		return nil, err
	}

	return it, nil
}

func (p *Plugin) showForm(l string, c *gin.Context) (interface{}, error) {
	var it = Form{}
	if err := p.DB.Where("id = ?", c.Param("id")).First(&it).Error; err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) getForm(c *gin.Context, id interface{}) (*Form, error) {
	user := c.MustGet(nut.CurrentUser).(*nut.User)
	var it = Form{}
	if err := p.DB.Where("id = ?", id).First(&it).Error; err != nil {
		return nil, err
	}
	if it.UserID != user.ID && !c.MustGet(nut.IsAdmin).(bool) {
		lang := c.MustGet(web.LOCALE).(string)
		return nil, p.I18n.E(lang, "errors.not-allow")
	}
	return &it, nil
}

func (p *Plugin) canEditForm(c *gin.Context) {
	it, err := p.getForm(c, c.Param("id"))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}
	c.Set("form", it)
}

func (p *Plugin) updateForm(l string, c *gin.Context) (interface{}, error) {
	var fm fmForm
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	it := c.MustGet("form").(*Form)

	begin, end, err := fm.parse()
	if err != nil {
		return nil, err
	}

	if err := p.DB.Model(it).Updates(map[string]interface{}{
		"title":     fm.Title,
		"type":      fm.Type,
		"body":      fm.Body,
		"mode":      fm.Mode,
		"start_up":  *begin,
		"shut_down": *end,
	}).Error; err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

func (p *Plugin) destroyForm(l string, c *gin.Context) (interface{}, error) {
	it := c.MustGet("form").(*Form)
	db := p.DB.Begin()
	if err := db.Where("form_id = ?", it.ID).Delete(Field{}).Error; err != nil {
		db.Rollback()
		return nil, err
	}
	if err := db.Where("form_id = ?", it.ID).Delete(Record{}).Error; err != nil {
		db.Rollback()
		return nil, err
	}
	if err := db.Delete(it).Error; err != nil {
		db.Rollback()
		return nil, err
	}
	db.Commit()
	return gin.H{}, nil
}
