package reading

import "github.com/gin-gonic/gin"

func (p *Plugin) indexBooks(l string, c *gin.Context) (interface{}, error) {
	var items []Book
	if err := p.DB.Select([]string{"id", "title"}).Order("updated_at DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
func (p *Plugin) destroyBook(l string, c *gin.Context) (interface{}, error) {
	bid := c.Param("id")

	db := p.DB.Begin()
	if err := db.Where("book_id = ?", bid).Delete(Note{}).Error; err != nil {
		db.Rollback()
		return nil, err
	}
	if err := db.Where("id = ?", bid).Delete(Book{}).Error; err != nil {
		db.Rollback()
		return nil, err
	}
	db.Commit()
	return gin.H{}, nil
}
