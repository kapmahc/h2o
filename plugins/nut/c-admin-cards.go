package nut

import "github.com/gin-gonic/gin"

func (p *Plugin) indexAdminCards(l string, c *gin.Context) (interface{}, error) {
	var items []Card
	if err := p.DB.Where("lang = ?", l).Order("loc ASC, sort_order ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

type fmCard struct {
	Href      string `json:"href" binding:"required"`
	Title     string `json:"title" binding:"required"`
	Summary   string `json:"summary" binding:"required"`
	Type      string `json:"type" binding:"required"`
	Action    string `json:"action" binding:"required"`
	Logo      string `json:"logo" binding:"required"`
	Loc       string `json:"loc" binding:"required"`
	SortOrder int    `json:"sortOrder"`
}

func (p *Plugin) createAdminCard(l string, c *gin.Context) (interface{}, error) {
	var fm fmCard
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	it := Card{
		Href:      fm.Href,
		Title:     fm.Title,
		Summary:   fm.Summary,
		Type:      fm.Type,
		Action:    fm.Action,
		Logo:      fm.Logo,
		Loc:       fm.Loc,
		SortOrder: fm.SortOrder,
		Lang:      l,
	}
	if err := p.DB.Create(&it).Error; err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) showAdminCard(l string, c *gin.Context) (interface{}, error) {
	var it = Card{}
	if err := p.DB.Where("id = ?", c.Param("id")).First(&it).Error; err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) updateAdminCard(l string, c *gin.Context) (interface{}, error) {
	var fm fmCard
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	var it = Card{}
	if err := p.DB.Where("id = ?", c.Param("id")).First(&it).Error; err != nil {
		return nil, err
	}
	if err := p.DB.Model(&it).Updates(map[string]interface{}{
		"title":      fm.Title,
		"type":       fm.Type,
		"action":     fm.Action,
		"summary":    fm.Summary,
		"logo":       fm.Logo,
		"href":       fm.Href,
		"loc":        fm.Loc,
		"sort_order": fm.SortOrder,
	}).Error; err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

func (p *Plugin) destroyAdminCard(l string, c *gin.Context) (interface{}, error) {
	if err := p.DB.Where("id = ?", c.Param("id")).Delete(Card{}).Error; err != nil {
		return nil, err
	}
	return gin.H{}, nil
}
