package survey

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (p *Plugin) indexRecords(l string, c *gin.Context) (interface{}, error) {
	form, err := p.getForm(c, c.Query("formId"))
	if err != nil {
		return nil, err
	}
	var items []Record
	if err := p.DB.Where("form_id = ?", form.ID).Order("updated_at DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (p *Plugin) canEditRecord(c *gin.Context) {
	var it = Record{}
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

	c.Set("record", &it)
}

func (p *Plugin) destroyRecord(l string, c *gin.Context) (interface{}, error) {
	it := c.MustGet("record").(*Record)
	if err := p.DB.Delete(it).Error; err != nil {
		return nil, err
	}

	return gin.H{}, nil
}
