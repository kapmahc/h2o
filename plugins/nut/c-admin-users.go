package nut

import "github.com/gin-gonic/gin"

func (p *Plugin) indexAdminUsers(l string, c *gin.Context) (interface{}, error) {
	var items []User
	if err := p.DB.Order("last_sign_in_at ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
