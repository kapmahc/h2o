package main

import (
	"os"

	_ "github.com/kapmahc/h2o/plugins/erp"
	_ "github.com/kapmahc/h2o/plugins/forum"
	_ "github.com/kapmahc/h2o/plugins/mall"
	_ "github.com/kapmahc/h2o/plugins/nut"
	_ "github.com/kapmahc/h2o/plugins/ops/mail"
	_ "github.com/kapmahc/h2o/plugins/ops/vpn"
	_ "github.com/kapmahc/h2o/plugins/pos"
	_ "github.com/kapmahc/h2o/plugins/reading"
	_ "github.com/kapmahc/h2o/plugins/survey"
	"github.com/kapmahc/h2o/web"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := web.Main(os.Args...); err != nil {
		log.Error(err)
	}
}
