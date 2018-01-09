package reading

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kapmahc/epub"
	"github.com/kapmahc/h2o/web"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// Shell console commands
func (p *Plugin) Shell() []cli.Command {
	return []cli.Command{
		{
			Name:    "books",
			Aliases: []string{"bk"},
			Usage:   "books operation",
			Subcommands: []cli.Command{
				{
					Name:    "sync",
					Aliases: []string{"s"},
					Usage:   fmt.Sprintf("sync books from %s", p.root()),
					Action: web.InjectAction(func(c *cli.Context) error {
						db := p.DB.Begin()
						if err := filepath.Walk(p.root(), func(path string, info os.FileInfo, err error) error {
							if err != nil {
								return err
							}
							if info.IsDir() {
								return nil
							}

							const epubExt = ".epub"
							const sep = ","

							name := info.Name()
							switch filepath.Ext(name) {
							case epubExt:
								log.Infof("find book %s", path)
								bk, err := epub.Open(path)
								if err != nil {
									return err
								}
								defer bk.Close()
								mt := bk.Opf.Metadata
								now := time.Now()
								var author []string
								for _, a := range mt.Creator {
									author = append(author, a.Data)
								}

								it := Book{
									File:        path[len(p.root())+1:],
									Type:        bk.Mimetype,
									Subject:     strings.Join(mt.Subject, sep),
									Author:      strings.Join(author, sep),
									Description: strings.Join(mt.Description, sep),
									Title:       strings.Join(mt.Title, sep),
									Lang:        strings.Join(mt.Language, sep),
									Publisher:   strings.Join(mt.Publisher, sep),
									UpdatedAt:   now,
									CreatedAt:   now,
								}
								if len(mt.Date) > 0 {
									it.PublishedAt, _ = time.Parse("2016-01-02", mt.Date[0].Data)
								}
								var cnt int
								if err := db.Model(&Book{}).Where("file = ?", it.File).Count(&cnt).Error; err != nil {
									db.Rollback()
									return err
								}
								if cnt > 0 {
									log.Warn("already exist!")
									return nil
								}
								if err := db.Create(&it).Error; err != nil {
									db.Rollback()
									return err
								}
							default:
								log.Warnf("ingnore file %s", name)
							}

							return nil
						}); err != nil {
							return err
						}
						db.Commit()
						// var items []Book
						// if err := p.DB.Select([]string{"uid", "email", "name"}).
						// 	Order("last_sign_in_at DESC").
						// 	Find(&items).Error; err != nil {
						// 	return err
						// }
						return nil
					}),
				},
			},
		},
	}
}
