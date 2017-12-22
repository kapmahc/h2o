package forum

import (
	"github.com/gin-gonic/gin"
	"github.com/kapmahc/h2o/plugins/nut"
)

func (p *Plugin) indexCommentsH(l string, c *gin.Context) (gin.H, error) {
	var items []Comment
	if err := p.DB.Order("updated_at DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	for id, it := range items {
		var a Article
		if err := p.DB.Select([]string{"id", "title"}).Where("id = ?", it.ArticleID).First(&a).Error; err != nil {
			return nil, err
		}
		items[id].Article = a
	}
	return gin.H{
		"comments": items,
		nut.TITLE:  p.I18n.T(l, "forum.comments.index.title"),
	}, nil
}

func (p *Plugin) showArticleH(l string, c *gin.Context) (gin.H, error) {
	var it Article
	if err := p.DB.Where("id = ?", c.Param("id")).First(&it).Error; err != nil {
		return nil, err
	}
	if err := p.DB.Model(&it).Association("Tags").Find(&it.Tags).Error; err != nil {
		return nil, err
	}
	if err := p.DB.Model(&it).Order("updated_at DESC").Related(&it.Comments).Error; err != nil {
		return nil, err
	}
	return gin.H{
		"article": it,
		"title":   it.Title,
	}, nil
}

func (p *Plugin) indexArticlesH(l string, c *gin.Context) (gin.H, error) {
	var items []Article
	if err := p.DB.Select([]string{"id", "title", "body"}).Order("updated_at DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return gin.H{
		"articles": items,
		"title":    p.I18n.T(l, "forum.articles.index.title"),
	}, nil
}

func (p *Plugin) showTagH(l string, c *gin.Context) (gin.H, error) {
	var it Tag
	if err := p.DB.Where("id = ?", c.Param("id")).Find(&it).Error; err != nil {
		return nil, err
	}
	if err := p.DB.Model(&it).Association("Articles").Find(&it.Articles).Error; err != nil {
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
