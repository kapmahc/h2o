package nut

import (
	"errors"
	"fmt"
	"html/template"
	"path"
	"time"

	"github.com/SermoDigital/jose/crypto"
	"github.com/facebookgo/inject"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/kapmahc/h2o/web"
	"github.com/spf13/viper"
	"github.com/unrolled/render"
)

func (p *Plugin) openDB() (*gorm.DB, error) {
	drv, url := p.database()
	db, err := gorm.Open(drv, url)
	if err != nil {
		return nil, err
	}
	if web.MODE() != web.PRODUCTION {
		db.LogMode(true)
	}
	return db, nil
}

func (p *Plugin) openS3() (*web.S3, error) {
	args := viper.GetStringMapString("aws")
	return web.NewS3(args["access_key_id"], args["secret_access_key"], args["region"], args["bucket"])
}

func (p *Plugin) openJobber() (*web.Jobber, error) {
	args := viper.GetStringMap("rabbitmq")
	return web.NewJobber(fmt.Sprintf(
		"amqp://%s:%s@%s:%d/%s",
		args["user"],
		args["password"],
		args["host"],
		args["port"],
		args["virtual"],
	), args["queue"].(string))
}

func (p *Plugin) openRender(theme string) *render.Render {
	helpers := template.FuncMap{
		"fmt": fmt.Sprintf,
		"dtf": func(t time.Time) string {
			return t.Format(time.RFC822)
		},
		"eq": func(a interface{}, b interface{}) bool {
			return a == b
		},
		"substr": func(s string, l int) string {
			return s[0:l]
		},
		"str2htm": func(s string) template.HTML {
			return template.HTML(s)
		},
		"site": func(k string) interface{} {
			switch k {
			case "version":
				return web.Version
			case "author":
				var author map[string]string
				if err := p.Settings.Get(p.DB, "site.author", &author); err != nil {
					author = map[string]string{
						"name":  "",
						"email": "",
					}
				}
				return author
			case "favicon":
				var favicon string
				if err := p.Settings.Get(p.DB, "site.favicon", &favicon); err != nil {
					favicon = "/assets/favicon.png"
				}
				return favicon
			case "languages":
				langs, _ := p.I18n.Languages(p.DB)
				return langs
			default:
				return k
			}
		},
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"t": func(lang, code string, args ...interface{}) string {
			return p.I18n.T(lang, code, args...)
		},
		"assets_css": func(u string) template.HTML {
			return template.HTML(fmt.Sprintf(`<link type="text/css" rel="stylesheet" href="%s">`, u))
		},
		"assets_js": func(u string) template.HTML {
			return template.HTML(fmt.Sprintf(`<script src="%s"></script>`, u))
		},
		"links": func(lng, loc string) ([]Link, error) {
			var items []Link
			if err := p.DB.Select([]string{"id", "label", "href", "loc", "sort_order"}).
				Where("lang = ? AND loc = ?", lng, loc).
				Order("sort_order ASC").
				Find(&items).Error; err != nil {
				return nil, err
			}
			return items, nil
		},
		"cards": func(lng, loc string) ([]Card, error) {
			var items []Card
			if err := p.DB.Select([]string{"id", "title", "summary", "type", "action", "logo", "href", "loc", "sort_order"}).
				Where("lang = ? AND loc = ?", lng, loc).
				Order("sort_order ASC").
				Find(&items).Error; err != nil {
				return nil, err
			}
			return items, nil
		},
		"odd": func(v int) bool {
			return v%2 != 0
		},
		"even": func(v int) bool {
			return v%2 == 0
		},
	}

	return render.New(render.Options{
		Directory:  path.Join("themes", theme, "views"),
		Extensions: []string{".html"},
		Layout:     "layouts/application/index",
		Funcs:      []template.FuncMap{helpers},
	})
}

func (p *Plugin) openRouter(db *gorm.DB, theme string) (*gin.Engine, error) {
	if web.MODE() == web.PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	}
	rt := gin.Default()
	i18m, err := p.I18n.Middleware(db)
	if err != nil {
		return nil, err
	}
	rt.Use(i18m)
	rt.Use(p.Layout.CurrentUserMiddleware)

	if web.MODE() != web.PRODUCTION {
		for k, v := range map[string]string{
			"3rd":    "node_modules",
			"assets": path.Join("themes", theme, "assets"),
		} {
			rt.Static("/"+k+"/", v)
		}
	}
	return rt, nil
}

func (p *Plugin) openRedis() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, e := redis.Dial(
				"tcp",
				fmt.Sprintf(
					"%s:%d",
					viper.GetString("redis.host"),
					viper.GetInt("redis.port"),
				),
			)
			if e != nil {
				return nil, e
			}
			if _, e = c.Do("SELECT", viper.GetInt("redis.db")); e != nil {
				c.Close()
				return nil, e
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

// Init init beans
func (p *Plugin) Init(g *inject.Graph) error {
	theme := viper.GetString("server.theme")

	db, err := p.openDB()
	if err != nil {
		return err
	}
	secret, err := web.SECRET()
	if err != nil {
		return err
	}

	security, err := web.NewSecurity(secret)
	if err != nil {
		return err
	}

	jobber, err := p.openJobber()
	if err != nil {
		return err
	}
	redis := p.openRedis()

	rt, err := p.openRouter(db, theme)
	if err != nil {
		return err
	}

	s3, err := p.openS3()
	if err != nil {
		return err
	}

	return g.Provide(
		&inject.Object{Value: db},
		&inject.Object{Value: redis},
		&inject.Object{Value: security},
		&inject.Object{Value: jobber},
		&inject.Object{Value: s3},
		&inject.Object{Value: web.NewCache("cache://")},
		&inject.Object{Value: web.NewJwt(secret, crypto.SigningMethodHS512)},
		&inject.Object{Value: rt},
		&inject.Object{Value: p.openRender(theme)},
		&inject.Object{Value: sessions.NewCookieStore(secret)},
	)
}
