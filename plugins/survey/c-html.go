package survey

import (
	"github.com/gin-gonic/gin"
	"github.com/kapmahc/h2o/plugins/nut"
)

func (p *Plugin) getApplyForm(l string, c *gin.Context) (gin.H, error) {
	it, err := p.parseForm(c.Param("id"))
	if err != nil {
		return nil, err
	}
	return gin.H{
		"form":    it,
		nut.TITLE: it.Title,
	}, nil
}

func (p *Plugin) getEditForm(l string, c *gin.Context) (gin.H, error) {
	it, err := p.parseForm(c.Param("id"))
	if err != nil {
		return nil, err
	}
	return gin.H{
		"form":    it,
		nut.TITLE: it.Title,
	}, nil
}

func (p *Plugin) getCancelForm(l string, c *gin.Context) error {
	return nil
}

func (p *Plugin) parseForm(id string) (*Form, error) {
	var it Form
	if err := p.DB.Where("id = ?", id).First(&it).Error; err != nil {
		return nil, err
	}
	if err := p.DB.Model(&it).Order("sort_order ASC").Related(&it.Fields).Error; err != nil {
		return nil, err
	}
	return &it, nil
}
