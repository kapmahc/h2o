package nut

import (
	"fmt"
	"time"

	"github.com/SermoDigital/jose/crypto"
	"github.com/facebookgo/inject"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/kapmahc/h2o/web"
	"github.com/spf13/viper"
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

	s3, err := p.openS3()
	if err != nil {
		return err
	}

	if web.MODE() == web.PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	}

	return g.Provide(
		&inject.Object{Value: viper.GetStringSlice("languages"), Name: "languages"},
		&inject.Object{Value: db},
		&inject.Object{Value: redis},
		&inject.Object{Value: security},
		&inject.Object{Value: jobber},
		&inject.Object{Value: s3},
		&inject.Object{Value: web.NewSitemap()},
		&inject.Object{Value: web.NewRSS()},
		&inject.Object{Value: web.NewCache("cache://")},
		&inject.Object{Value: web.NewJwt(secret, crypto.SigningMethodHS512)},
		&inject.Object{Value: gin.Default()},
	)
}
