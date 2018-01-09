package reading

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kapmahc/epub"
	"github.com/kapmahc/h2o/plugins/nut"
	"github.com/kapmahc/h2o/web"
)

// http://www.cbeta.org/cbreader/help/cbr_toc.htm

func (p *Plugin) showBookH(l string, c *gin.Context) (gin.H, error) {
	var buf bytes.Buffer
	it, bk, err := p.readBook(c.Param("id"))
	if err != nil {
		return nil, err
	}
	if err := p.DB.Model(&it).Order("updated_at DESC").Related(&it.Notes).Error; err != nil {
		return nil, err
	}
	// c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")

	if len(bk.Ncx.Points) > 0 {
		p.writePoints(
			&buf,
			fmt.Sprintf("/reading/htdocs/pages/%d", it.ID),
			bk.Ncx.Points,
		)
	} else {
		p.writeManifest(
			&buf,
			fmt.Sprintf("/reading/htdocs/pages/%d", it.ID),
			bk.Opf.Manifest,
		)
	}
	return gin.H{
		"homeage": buf.String(),
		"book":    it,
		nut.TITLE: it.Title,
	}, nil
}

func (p *Plugin) indexBooksH(l string, c *gin.Context) (gin.H, error) {
	var total int64
	if err := p.DB.Model(&Book{}).Count(&total).Error; err != nil {
		return nil, err
	}

	pag := web.NewPagination(c, total)
	var items []Book

	if err := p.DB.
		Select([]string{"id", "title", "author", "subject", "description"}).
		Order("updated_at DESC").
		Offset((pag.Page - 1) * pag.Size).
		Limit(pag.Size).
		Find(&items).Error; err != nil {
		return nil, err
	}
	return gin.H{
		"pagination": pag,
		"books":      items,
		"total":      total,
		"title":      p.I18n.T(l, "reading.books.index.title"),
	}, nil
}

func (p *Plugin) showPage(c *gin.Context) {
	if err := p.readBookPage(c.Writer, c.Param("id"), c.Param("href")[1:]); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}

// -------------------

func (p *Plugin) readBook(id string) (*Book, *epub.Book, error) {
	var book Book
	if err := p.DB.
		Where("id = ?", id).First(&book).Error; err != nil {
		return nil, nil, err
	}
	bk, err := epub.Open(path.Join(p.root(), book.File))
	return &book, bk, err
}

func (p *Plugin) readBookPage(w http.ResponseWriter, id string, name string) error {
	_, bk, err := p.readBook(id)
	if err != nil {
		return err
	}
	for _, fn := range bk.Files() {
		if strings.HasSuffix(fn, name) {
			for _, mf := range bk.Opf.Manifest {
				log.Printf("%s %s", mf.Href, name)
				if mf.Href == name {
					rdr, err := bk.Open(name)
					if err != nil {
						return err
					}
					defer rdr.Close()
					body, err := ioutil.ReadAll(rdr)
					if err != nil {
						return err
					}
					w.Header().Set("Content-Type", mf.MediaType)
					w.Write(body)
					return nil
				}
			}
		}
	}
	return errors.New("not found")
}

func (p *Plugin) writeManifest(wrt io.Writer, href string, manifests []epub.Manifest) {
	wrt.Write([]byte("<ol>"))
	for _, it := range manifests {
		if it.MediaType == "application/xhtml+xml" {
			wrt.Write([]byte("<li>"))
			fmt.Fprintf(
				wrt,
				`<a href="%s/%s" target="_blank">%s</a>`,
				href,
				it.Href,
				it.Href,
			)
			wrt.Write([]byte("</li>"))
		}
	}
	wrt.Write([]byte("</ol>"))
}

func (p *Plugin) writePoints(wrt io.Writer, href string, points []epub.NavPoint) {
	wrt.Write([]byte("<ol>"))
	for _, it := range points {
		wrt.Write([]byte("<li>"))
		fmt.Fprintf(
			wrt,
			`<a href="%s/%s" target="_blank">%s</a>`,
			href,
			it.Content.Src,
			it.Text,
		)
		p.writePoints(wrt, href, it.Points)
		wrt.Write([]byte("</li>"))
	}
	wrt.Write([]byte("</ol>"))
}
