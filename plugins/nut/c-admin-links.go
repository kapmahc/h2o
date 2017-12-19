package nut

import "github.com/gin-gonic/gin"

func (p *Plugin) indexAdminLinks(l string, c *gin.Context) (interface{}, error) {
	var items []Link
	if err := p.DB.Where("lang = ?", l).Order("loc ASC, sort_order ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

type fmLink struct {
	Href      string `json:"href" binding:"required"`
	Label     string `json:"label" binding:"required"`
	Loc       string `json:"loc" binding:"required"`
	SortOrder int    `json:"sortOrder"`
}

func (p *Plugin) createAdminLink(l string, c *gin.Context) (interface{}, error) {
	var fm fmLink
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	it := Link{
		Href:      fm.Href,
		Label:     fm.Label,
		Loc:       fm.Loc,
		SortOrder: fm.SortOrder,
		Lang:      l,
	}
	if err := p.DB.Create(&it).Error; err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) showAdminLink(l string, c *gin.Context) (interface{}, error) {
	var it = Link{}
	if err := p.DB.Where("id = ?", c.Param("id")).First(&it).Error; err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) updateAdminLink(l string, c *gin.Context) (interface{}, error) {
	var fm fmLink
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	var it = Link{}
	if err := p.DB.Where("id = ?", c.Param("id")).First(&it).Error; err != nil {
		return nil, err
	}
	if err := p.DB.Model(&it).Updates(map[string]interface{}{
		"label":      fm.Label,
		"href":       fm.Href,
		"loc":        fm.Loc,
		"sort_order": fm.SortOrder,
	}).Error; err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

func (p *Plugin) destroyAdminLink(l string, c *gin.Context) (interface{}, error) {
	if err := p.DB.Where("id = ?", c.Param("id")).Delete(Link{}).Error; err != nil {
		return nil, err
	}
	return gin.H{}, nil
}
