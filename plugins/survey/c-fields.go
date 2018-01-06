package survey

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (p *Plugin) indexFields(l string, c *gin.Context) (interface{}, error) {
	var items []Field
	if err := p.DB.Where("form_id = ?", c.Query("formId")).Order("sort_order ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

type fmField struct {
	Label     string `json:"label" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Type      string `json:"type" binding:"required"`
	Value     string `json:"value"`
	Body      string `json:"body"`
	Required  bool   `json:"required"`
	SortOrder int    `json:"sortOrder"`
}

func (p *fmField) parse() ([]byte, []byte, error) {
	val, err := json.Marshal(strings.Split(p.Value, ";"))
	if err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(strings.Split(p.Body, "\n"))
	if err != nil {
		return nil, nil, err
	}
	return val, body, nil
}

func (p *Plugin) createField(l string, c *gin.Context) (interface{}, error) {
	var fm fmField
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	val, bod, err := fm.parse()
	if err != nil {
		return nil, err
	}
	form, err := p.getForm(c, c.Query("formId"))
	if err != nil {
		return nil, err
	}

	it := Field{
		Label:     fm.Label,
		Name:      fm.Name,
		Body:      string(bod),
		Type:      fm.Type,
		Value:     string(val),
		Required:  fm.Required,
		SortOrder: fm.SortOrder,
		FormID:    form.ID,
	}

	if err := p.DB.Create(&it).Error; err != nil {
		return nil, err
	}

	return it, nil
}

func (p *Plugin) showField(l string, c *gin.Context) (interface{}, error) {
	var it = Field{}
	if err := p.DB.Where("id = ?", c.Param("id")).First(&it).Error; err != nil {
		return nil, err
	}
	return it, nil
}

func (p *Plugin) canEditField(c *gin.Context) {
	var it = Field{}
	if err := p.DB.Where("id = ?", c.Param("id")).First(&it).Error; err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}
	_, err := p.getForm(c, it.FormID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}

	c.Set("field", &it)
}

func (p *Plugin) updateField(l string, c *gin.Context) (interface{}, error) {
	var fm fmField
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	val, bod, err := fm.parse()
	if err != nil {
		return nil, err
	}
	it := c.MustGet("field").(*Field)

	if err := p.DB.Model(it).Updates(map[string]interface{}{
		"label":      fm.Label,
		"name":       fm.Name,
		"type":       fm.Type,
		"body":       string(bod),
		"value":      string(val),
		"required":   fm.Required,
		"sort_order": fm.SortOrder,
	}).Error; err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

func (p *Plugin) destroyField(l string, c *gin.Context) (interface{}, error) {
	it := c.MustGet("field").(*Field)
	if err := p.DB.Delete(it).Error; err != nil {
		return nil, err
	}

	return gin.H{}, nil
}
