package nut

import (
	"github.com/gin-gonic/gin"
	"github.com/kapmahc/h2o/web"
)

func (p *Plugin) indexAdminLocales(l string, c *gin.Context) (interface{}, error) {
	var items []web.Locale
	if err := p.DB.Select([]string{"id", "code", "message"}).Where("lang = ?", l).Order("code ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

type fmLocale struct {
	Code    string `json:"code" binding:"required"`
	Message string `json:"message" binding:"required"`
}

func (p *Plugin) createAdminLocale(l string, c *gin.Context) (interface{}, error) {
	var fm fmLocale
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	if err := p.I18n.Set(p.DB, l, fm.Code, fm.Message); err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

func (p *Plugin) showAdminLocale(l string, c *gin.Context) (interface{}, error) {
	var it web.Locale
	if err := p.DB.Where("id = ?", c.Param("id")).First(&it).Error; err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) destroyAdminLocale(l string, c *gin.Context) (interface{}, error) {
	if err := p.DB.Where("id = ?", c.Param("id")).Delete(web.Locale{}).Error; err != nil {
		return nil, err
	}
	return gin.H{}, nil
}
