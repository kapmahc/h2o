#!/bin/sh
go get -u github.com/kardianos/govendor
rm -r vendor
govendor init
govendor fetch github.com/unrolled/render
govendor fetch github.com/gorilla/sessions
govendor fetch github.com/gorilla/csrf
govendor fetch github.com/gin-gonic/gin
govendor fetch github.com/jinzhu/gorm
govendor fetch github.com/jinzhu/gorm/dialects/postgres
govendor fetch github.com/jinzhu/gorm/dialects/mysql
govendor fetch github.com/mattes/migrate
govendor fetch github.com/mattes/migrate/database/postgres
govendor fetch github.com/mattes/migrate/database/mysql
govendor fetch github.com/mattes/migrate/source/file
govendor fetch golang.org/x/crypto/bcrypt
govendor fetch golang.org/x/text/language
govendor fetch github.com/SermoDigital/jose/jwt
govendor fetch github.com/SermoDigital/jose/jws
govendor fetch github.com/SermoDigital/jose/crypto
govendor fetch github.com/streadway/amqp
govendor fetch github.com/google/uuid
govendor fetch github.com/spf13/viper
govendor fetch github.com/urfave/cli
govendor fetch github.com/facebookgo/inject
govendor fetch github.com/garyburd/redigo/redis
govendor fetch github.com/go-ini/ini
govendor fetch github.com/sirupsen/logrus
govendor fetch github.com/sirupsen/logrus/hooks/syslog
govendor fetch github.com/BurntSushi/toml
govendor fetch gopkg.in/gomail.v2
govendor fetch github.com/aws/aws-sdk-go/aws/session
govendor fetch github.com/aws/aws-sdk-go/service/s3
govendor fetch github.com/kapmahc/stardict
govendor fetch github.com/kapmahc/epub
