package nut

import (
	"encoding/base64"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/kapmahc/h2o/web"
	"github.com/spf13/viper"
	"github.com/unrolled/render"
	"golang.org/x/text/language"
)

// Plugin plugin
type Plugin struct {
	I18n     *web.I18n      `inject:""`
	Cache    *web.Cache     `inject:""`
	Jwt      *web.Jwt       `inject:""`
	Jobber   *web.Jobber    `inject:""`
	Settings *web.Settings  `inject:""`
	Router   *gin.Engine    `inject:""`
	DB       *gorm.DB       `inject:""`
	Render   *render.Render `inject:""`
	Dao      *Dao           `inject:""`
	Layout   *Layout        `inject:""`
}

func init() {
	web.Register(&Plugin{})

	viper.SetDefault("languages", []string{
		language.AmericanEnglish.String(),
		language.SimplifiedChinese.String(),
		language.TraditionalChinese.String(),
	})

	viper.SetDefault("aws", map[string]interface{}{
		"access_key_id":     "change-me",
		"secret_access_key": "change-me",
		"region":            "change-me",
		"bucket":            "change-me",
	})

	viper.SetDefault("redis", map[string]interface{}{
		"host": "localhost",
		"port": 6379,
		"db":   8,
	})

	viper.SetDefault("rabbitmq", map[string]interface{}{
		"user":     "guest",
		"password": "guest",
		"host":     "localhost",
		"port":     5672,
		"virtual":  "h2o-dev",
		"queue":    "tasks",
	})

	viper.SetDefault("database", map[string]interface{}{
		"driver": "postgres",
		"args": map[string]interface{}{
			"host":     "localhost",
			"port":     5432,
			"user":     "postgres",
			"password": "",
			"dbname":   "h2o_dev",
			"sslmode":  "disable",
		},
	})

	viper.SetDefault("server", map[string]interface{}{
		"port":     8080,
		"frontend": "http://localhost:3000",
		"backend":  "http://localhost:8080",
		"name":     "change-me.com",
		"theme":    "bootstrap",
		"secure":   false,
	})

	secret, _ := web.RandomBytes(32)
	viper.SetDefault("secret", base64.StdEncoding.EncodeToString(secret))

	viper.SetDefault("elasticsearch", []string{"http://localhost:9200"})
}
