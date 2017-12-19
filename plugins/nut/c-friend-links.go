package nut

import "github.com/gin-gonic/gin"

func (p *Plugin) indexAdminFriendLinks(l string, c *gin.Context) (interface{}, error) {
	var items []FriendLink
	if err := p.DB.Order("sort_order ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

type fmFriendLink struct {
	Title     string `json:"title" binding:"required"`
	Home      string `json:"home" binding:"required"`
	Logo      string `json:"logo" binding:"required"`
	SortOrder int    `json:"sortOrder"`
}

func (p *Plugin) createAdminFriendLink(l string, c *gin.Context) (interface{}, error) {
	var fm fmFriendLink
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	it := FriendLink{
		Title:     fm.Title,
		Home:      fm.Home,
		Logo:      fm.Logo,
		SortOrder: fm.SortOrder,
	}
	if err := p.DB.Create(&it).Error; err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) showAdminFriendLink(l string, c *gin.Context) (interface{}, error) {
	var it = FriendLink{}
	if err := p.DB.Where("id = ?", c.Param("id")).First(&it).Error; err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) updateAdminFriendLink(l string, c *gin.Context) (interface{}, error) {
	var fm fmFriendLink
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	var it = FriendLink{}
	if err := p.DB.Where("id = ?", c.Param("id")).First(&it).Error; err != nil {
		return nil, err
	}
	if err := p.DB.Model(&it).Updates(map[string]interface{}{
		"title":      fm.Title,
		"home":       fm.Home,
		"logo":       fm.Logo,
		"sort_order": fm.SortOrder,
	}).Error; err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

func (p *Plugin) destroyAdminFriendLink(l string, c *gin.Context) (interface{}, error) {
	if err := p.DB.Where("id = ?", c.Param("id")).Delete(FriendLink{}).Error; err != nil {
		return nil, err
	}
	return gin.H{}, nil
}
