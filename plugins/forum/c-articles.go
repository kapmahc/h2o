package forum

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kapmahc/h2o/plugins/nut"
	"github.com/kapmahc/h2o/web"
)

func (p *Plugin) latestArticles(l string, c *gin.Context) (interface{}, error) {
	var items []Article
	if err := p.DB.Select([]string{"title", "id", "body"}).
		Order("updated_at DESC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (p *Plugin) indexArticles(l string, c *gin.Context) (interface{}, error) {
	var items []Article
	db := p.DB
	user := c.MustGet(nut.CurrentUser).(*nut.User)
	if c.MustGet(nut.IsAdmin).(bool) {
		db = p.DB.Where("lang = ?", l)
	} else {
		p.DB.Where("lang = ? AND user_id = ?", l, user.ID)
	}
	if err := db.Order("updated_at DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

type fmArticle struct {
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
	Type  string `json:"type" binding:"required"`
	Tags  []uint `json:"tags"`
}

func (p *fmArticle) toTags() []interface{} {
	var tags []interface{}
	for _, id := range p.Tags {
		tags = append(tags, Tag{ID: id})
	}
	return tags
}

func (p *Plugin) createArticle(l string, c *gin.Context) (interface{}, error) {
	user := c.MustGet(nut.CurrentUser).(*nut.User)
	var fm fmArticle
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	it := Article{
		Title:  fm.Title,
		Body:   fm.Body,
		Type:   fm.Type,
		Lang:   l,
		UserID: user.ID,
	}
	db := p.DB.Begin()
	if err := db.Create(&it).Error; err != nil {
		db.Rollback()
		return nil, err
	}

	tags := fm.toTags()
	if len(tags) > 0 {
		if err := db.Model(&it).Association("Tags").Append(tags...).Error; err != nil {
			db.Rollback()
			return nil, err
		}
	}

	db.Commit()
	return it, nil
}

func (p *Plugin) showArticle(l string, c *gin.Context) (interface{}, error) {
	var it = Article{}
	if err := p.DB.Where("id = ?", c.Param("id")).First(&it).Error; err != nil {
		return nil, err
	}
	if err := p.DB.Model(it).Association("Tags").Find(&it.Tags).Error; err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) canEditArticle(c *gin.Context) {
	user := c.MustGet(nut.CurrentUser).(*nut.User)
	var it = Article{}
	if err := p.DB.Where("id = ?", c.Param("id")).First(&it).Error; err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}
	lang := c.MustGet(web.LOCALE).(string)
	if it.UserID != user.ID && !c.MustGet(nut.IsAdmin).(bool) {
		c.String(http.StatusForbidden, p.I18n.T(lang, "errors.not-allow"))
		c.Abort()
		return
	}
	c.Set("article", &it)
}

func (p *Plugin) updateArticle(l string, c *gin.Context) (interface{}, error) {
	var fm fmArticle
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	it := c.MustGet("article").(*Article)

	db := p.DB.Begin()
	if err := db.Model(it).Updates(map[string]interface{}{
		"title": fm.Title,
		"type":  fm.Type,
		"body":  fm.Body,
	}).Error; err != nil {
		db.Rollback()
		return nil, err
	}
	tags := fm.toTags()
	if len(tags) > 0 {
		if err := db.Model(&it).Association("Tags").Replace(tags...).Error; err != nil {
			db.Rollback()
			return nil, err
		}
	}
	db.Commit()
	return gin.H{}, nil
}

func (p *Plugin) destroyArticle(l string, c *gin.Context) (interface{}, error) {
	it := c.MustGet("article").(*Article)
	db := p.DB.Begin()
	if err := db.Where("article_id = ?", it.ID).Delete(Comment{}).Error; err != nil {
		db.Rollback()
		return nil, err
	}
	if err := db.Model(it).Association("Tags").Clear().Error; err != nil {
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
