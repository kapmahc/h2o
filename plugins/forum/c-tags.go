package forum

import "github.com/gin-gonic/gin"

func (p *Plugin) showTagH(l string, c *gin.Context) (gin.H, error) {
	var it Tag
	if err := p.DB.Where("id = ?", c.Param("id")).Find(&it).Error; err != nil {
		return nil, err
	}
	return gin.H{
		"tag":   it,
		"title": it.Name,
	}, nil
}

func (p *Plugin) indexTagsH(l string, c *gin.Context) (gin.H, error) {
	var items []Tag
	if err := p.DB.Select([]string{"id", "name"}).Order("updated_at DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return gin.H{
		"tags":  items,
		"title": p.I18n.T(l, "forum.tags.index.title"),
	}, nil
}

func (p *Plugin) indexTags(l string, c *gin.Context) (interface{}, error) {
	var items []Tag
	if err := p.DB.Order("updated_at DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

type fmTag struct {
	Name string `json:"name" binding:"required"`
}

func (p *Plugin) createTag(l string, c *gin.Context) (interface{}, error) {
	var fm fmTag
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	it := Tag{
		Name: fm.Name,
	}
	if err := p.DB.Create(&it).Error; err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) showTag(l string, c *gin.Context) (interface{}, error) {
	var it = Tag{}
	if err := p.DB.Where("id = ?", c.Param("id")).First(&it).Error; err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) updateTag(l string, c *gin.Context) (interface{}, error) {
	var fm fmTag
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	var it = Tag{}
	if err := p.DB.Where("id = ?", c.Param("id")).First(&it).Error; err != nil {
		return nil, err
	}
	if err := p.DB.Model(&it).Update("name", fm.Name).Error; err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

func (p *Plugin) destroyTag(l string, c *gin.Context) (interface{}, error) {
	db := p.DB.Begin()
	var it Tag
	if err := db.Where("id = ?", c.Param("id")).First(&it).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&it).Association("Articles").Clear().Error; err != nil {
		db.Rollback()
		return nil, err
	}
	if err := db.Delete(&it).Error; err != nil {
		db.Rollback()
		return nil, err
	}
	db.Commit()
	return gin.H{}, nil
}
