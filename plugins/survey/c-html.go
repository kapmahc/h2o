package survey

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kapmahc/h2o/plugins/nut"
)

func (p *Plugin) postApplyForm(l string, c *gin.Context) (interface{}, error) {
	it, _, err := p.selectForm(c.Param("id"))
	if err != nil {
		return nil, err
	}
	email, value, err := p.parseFormData(c.Request, it)
	if err != nil {
		return nil, err
	}
	if err := p.DB.Create(&Record{FormID: it.ID, Email: email, Value: string(value)}).Error; err != nil {
		return nil, err
	}
	return gin.H{"ok": true, "message": p.I18n.T(l, "helpers.success")}, nil
}

func (p *Plugin) parseFormData(r *http.Request, f *Form) (string, []byte, error) {
	err := r.ParseForm()
	if err != nil {
		return "", nil, err
	}
	val := gin.H{}
	for _, fd := range f.Fields {
		switch {
		case fd.Type == "text" || fd.Type == "textarea" || fd.Type == "select" || fd.Type == "radios":
			val[fd.Name] = r.Form.Get(fd.Name)
		case fd.Type == "checkboxes":
			val[fd.Name] = r.Form[fd.Name]
		}
	}
	buf, err := json.Marshal(val)
	if err != nil {
		return "", nil, err
	}
	return r.Form.Get("email"), buf, err
}

func (p *Plugin) getApplyForm(l string, c *gin.Context) (gin.H, error) {
	it, options, err := p.selectForm(c.Param("id"))
	if err != nil {
		return nil, err
	}
	return gin.H{
		"form":    it,
		"options": options,
		"values":  make(map[string]string),
		nut.TITLE: it.Title,
	}, nil
}

func (p *Plugin) getEditForm(l string, c *gin.Context) (gin.H, error) {
	it, options, err := p.selectForm(c.Param("id"))
	if err != nil {
		return nil, err
	}
	return gin.H{
		"form":    it,
		"options": options,
		nut.TITLE: it.Title,
	}, nil
}

func (p *Plugin) getCancelForm(l string, c *gin.Context) error {
	return nil
}

func (p *Plugin) selectForm(id string) (*Form, map[string][]string, error) {
	var it Form
	if err := p.DB.Where("id = ?", id).First(&it).Error; err != nil {
		return nil, nil, err
	}
	if err := p.DB.Model(&it).Order("sort_order ASC").Related(&it.Fields).Error; err != nil {
		return nil, nil, err
	}

	options := make(map[string][]string)
	for _, fd := range it.Fields {
		switch {
		case fd.Type == "select" || fd.Type == "checkboxes" || fd.Type == "radios":
			var val []string
			if err := json.Unmarshal([]byte(fd.Body), &val); err != nil {
				return nil, nil, err
			}
			options[fd.Name] = val
		}
	}
	return &it, options, nil
}
