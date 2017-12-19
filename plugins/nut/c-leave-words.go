package nut

import "github.com/gin-gonic/gin"

func (p *Plugin) createLeaveWord(l string, c *gin.Context) (interface{}, error) {
	var fm fmLeaveWord
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	it := LeaveWord{
		Body: fm.Body,
		Type: fm.Type,
	}
	if err := p.DB.Create(&it).Error; err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) indexAdminLeaveWords(l string, c *gin.Context) (interface{}, error) {
	var items []LeaveWord
	if err := p.DB.Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

type fmLeaveWord struct {
	Body string `json:"body" binding:"required"`
	Type string `json:"type" binding:"required"`
}

func (p *Plugin) destroyAdminLeaveWord(l string, c *gin.Context) (interface{}, error) {
	if err := p.DB.Where("id = ?", c.Param("id")).Delete(LeaveWord{}).Error; err != nil {
		return nil, err
	}
	return gin.H{}, nil
}
