package nut

import (
	"context"
	"crypto/x509/pkix"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/csrf"
	"github.com/kapmahc/h2o/web"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database"
	"github.com/mattes/migrate/database/postgres"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
	"golang.org/x/text/language"
)

const (
	postgresqlDriverName = "postgres"
)

// Shell console commands
func (p *Plugin) Shell() []cli.Command {
	return []cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "generate file template",
			Subcommands: []cli.Command{
				{
					Name:    "config",
					Aliases: []string{"c"},
					Usage:   "generate config file",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "environment, e",
							Value: "development",
							Usage: "environment, like: development, production, stage, test...",
						},
					},
					Action: p.generateConfig,
				},
				{
					Name:    "nginx",
					Aliases: []string{"ng"},
					Usage:   "generate nginx.conf",
					Action:  web.ConfigAction(p.generateNginxConf),
				},
				{
					Name:    "openssl",
					Aliases: []string{"ssl"},
					Usage:   "generate ssl certificates",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Usage: "name",
						},
						cli.StringFlag{
							Name:  "country, c",
							Value: "Earth",
							Usage: "country",
						},
						cli.StringFlag{
							Name:  "organization, o",
							Value: "Mother Nature",
							Usage: "organization",
						},
						cli.IntFlag{
							Name:  "years, y",
							Value: 1,
							Usage: "years",
						},
					},
					Action: p.generateSsl,
				},
				{
					Name:    "migration",
					Usage:   "generate migration file",
					Aliases: []string{"m"},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Usage: "name",
						},
					},
					Action: p.generateMigration,
				},
				{
					Name:    "locale",
					Usage:   "generate locale file",
					Aliases: []string{"l"},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Usage: "locale name",
						},
					},
					Action: p.generateLocale,
				},
			},
		},
		{
			Name:    "cache",
			Aliases: []string{"c"},
			Usage:   "cache operations",
			Subcommands: []cli.Command{
				{
					Name:    "list",
					Usage:   "list all cache keys",
					Aliases: []string{"l"},
					Action: web.InjectAction(func(_ *cli.Context) error {
						keys, err := p.Cache.Keys()
						if err != nil {
							return err
						}
						for _, k := range keys {
							fmt.Println(k)
						}
						return nil
					}),
				},
				{
					Name:    "clear",
					Usage:   "clear cache items",
					Aliases: []string{"c"},
					Action: web.InjectAction(func(_ *cli.Context) error {
						return p.Cache.Flush()
					}),
				},
			},
		},
		{
			Name:    "database",
			Aliases: []string{"db"},
			Usage:   "database operations",
			Subcommands: []cli.Command{
				{
					Name:    "example",
					Usage:   "scripts example for create database and user",
					Aliases: []string{"e"},
					Action:  web.ConfigAction(p.databaseExample),
				},
				{
					Name:    "migrate",
					Usage:   "migrate the DB to the most recent version available",
					Aliases: []string{"m"},
					Action: p.databaseRun(func(m *migrate.Migrate) error {
						return m.Up()
					}),
				},
				{
					Name:    "rollback",
					Usage:   "roll back the version by 1",
					Aliases: []string{"r"},
					Action: p.databaseRun(func(m *migrate.Migrate) error {
						return m.Steps(-1)
					}),
				},
				{
					Name:    "version",
					Usage:   "dump the migration status for the current DB",
					Aliases: []string{"v"},
					Action: p.databaseRun(func(m *migrate.Migrate) error {
						version, dirty, err := m.Version()
						if err != nil {
							return err
						}
						fmt.Print(version)
						if dirty {
							fmt.Print("dirty")
						}
						fmt.Println()
						return nil
					}),
				},
				{
					Name:    "connect",
					Usage:   "connect database",
					Aliases: []string{"c"},
					Action:  web.ConfigAction(p.connectDatabase),
				},
				{
					Name:    "create",
					Usage:   "create database",
					Aliases: []string{"n"},
					Action:  web.ConfigAction(p.createDatabase),
				},
				{
					Name:    "drop",
					Usage:   "drop database",
					Aliases: []string{"d"},
					Action:  web.ConfigAction(p.dropDatabase),
				},
			},
		},
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "start the app server",
			Action: web.InjectAction(func(_ *cli.Context) error {
				go func() {
					// ----------
					host, err := os.Hostname()
					if err != nil {
						log.Error(err)
					}
					for {
						if err := p.Jobber.Receive(host); err != nil {
							log.Error(err)
							time.Sleep(5 * time.Second)
						}
					}
				}()
				// -------
				return p.startServer(viper.GetInt("server.port"), viper.GetString("env") == web.PRODUCTION)
			}),
		},
		{
			Name:    "routes",
			Aliases: []string{"rt"},
			Usage:   "print out all defined routes",
			Action: web.InjectAction(func(_ *cli.Context) error {
				tpl := "%-16s %s\n"
				fmt.Printf(tpl, "METHODS", "PATH")
				for _, rt := range p.Router.Routes() {
					fmt.Printf(tpl, rt.Method, rt.Path)
				}
				return nil
			}),
		},
	}
}

// --------------------------------------------

func (p *Plugin) generateNginxConf(c *cli.Context) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	name := viper.GetString("server.name")

	fn := path.Join("tmp", "etc", "nginx", "sites-enabled", name+".conf")
	if err = os.MkdirAll(path.Dir(fn), 0700); err != nil {
		return err
	}
	fmt.Printf("generate file %s\n", fn)
	fd, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer fd.Close()

	tpl, err := template.ParseFiles(path.Join("templates", "nginx.conf"))
	if err != nil {
		return err
	}

	return tpl.Execute(fd, struct {
		Port  int
		Root  string
		Name  string
		Theme string
		Ssl   bool
	}{
		Name:  name,
		Port:  viper.GetInt("server.port"),
		Root:  pwd,
		Theme: viper.GetString("server.theme"),
		Ssl:   viper.GetBool("server.secure"),
	})
}
func (p *Plugin) generateSsl(c *cli.Context) error {
	name := c.String("name")
	if len(name) == 0 {
		cli.ShowCommandHelp(c, "openssl")
		return nil
	}
	root := path.Join("tmp", "etc", "ssl", name)

	key, crt, err := CreateCertificate(
		true,
		pkix.Name{
			Country:      []string{c.String("country")},
			Organization: []string{c.String("organization")},
		},
		c.Int("years"),
	)
	if err != nil {
		return err
	}

	fnk := path.Join(root, "key.pem")
	fnc := path.Join(root, "crt.pem")

	fmt.Printf("generate pem file %s\n", fnk)
	err = WritePemFile(fnk, "RSA PRIVATE KEY", key, 0600)
	fmt.Printf("test: openssl rsa -noout -text -in %s\n", fnk)

	if err == nil {
		fmt.Printf("generate pem file %s\n", fnc)
		err = WritePemFile(fnc, "CERTIFICATE", crt, 0444)
		fmt.Printf("test: openssl x509 -noout -text -in %s\n", fnc)
	}
	if err == nil {
		fmt.Printf(
			"verify: diff <(openssl rsa -noout -modulus -in %s) <(openssl x509 -noout -modulus -in %s)",
			fnk,
			fnc,
		)
	}
	fmt.Println()
	return nil
}
func (p *Plugin) generateLocale(c *cli.Context) error {
	name := c.String("name")
	if len(name) == 0 {
		cli.ShowCommandHelp(c, "locale")
		return nil
	}
	lng, err := language.Parse(name)
	if err != nil {
		return err
	}
	const root = "locales"
	if err = os.MkdirAll(root, 0700); err != nil {
		return err
	}
	file := path.Join(root, fmt.Sprintf("%s.ini", lng.String()))
	fmt.Printf("generate file %s\n", file)
	fd, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer fd.Close()
	return err
}
func (p *Plugin) migrationsDir() string {
	return filepath.Join("db", viper.GetString("database.driver"), "migrations")
}
func (p *Plugin) generateMigration(c *cli.Context) error {
	name := c.String("name")
	if len(name) == 0 {
		cli.ShowCommandHelp(c, "migration")
		return nil
	}
	root := p.migrationsDir()
	version := time.Now().Format("20060102150405")
	if err := os.MkdirAll(root, 0700); err != nil {
		return err
	}
	for _, v := range []string{"up", "down"} {
		fn := filepath.Join(root, fmt.Sprintf("%s_%s.%s.sql", version, name, v))
		fmt.Printf("generate file %s\n", fn)
		fd, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if err != nil {
			return err
		}
		defer fd.Close()
	}
	return nil
}
func (p *Plugin) generateConfig(c *cli.Context) error {
	const fn = "config.toml"
	if _, err := os.Stat(fn); err == nil {
		return fmt.Errorf("file %s already exists", fn)
	}
	fmt.Printf("generate file %s\n", fn)

	viper.Set("env", c.String("environment"))

	fd, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer fd.Close()

	enc := toml.NewEncoder(fd)
	return enc.Encode(viper.AllSettings())
}

// --------------------------------
func (p *Plugin) databaseExample(_ *cli.Context) error {
	drv := viper.GetString("database.driver")
	switch drv {
	case postgresqlDriverName:
		args := viper.GetStringMapString("database.args")
		fmt.Printf("CREATE USER %s WITH PASSWORD '%s';\n", args["user"], args["password"])
		fmt.Printf("CREATE DATABASE %s WITH ENCODING='UTF8';\n", args["dbname"])
		fmt.Printf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s;\n", args["dbname"], args["user"])
		return nil
	default:
		return fmt.Errorf("bad driver %s", drv)
	}
}

func (p *Plugin) database() (string, string) {
	drv := viper.GetString("database.driver")
	args := viper.GetStringMap("database.args")
	url := ""
	for k, v := range args {
		url += fmt.Sprintf("%s=%v ", k, v)
	}
	return drv, url
}

func (p *Plugin) databaseRun(f func(*migrate.Migrate) error) cli.ActionFunc {
	return web.ConfigAction(func(_ *cli.Context) error {
		drv, url := p.database()
		db, err := sql.Open(drv, url)
		if err != nil {
			return err
		}
		defer db.Close()

		var ins database.Driver
		switch drv {
		case postgresqlDriverName:
			ins, err = postgres.WithInstance(db, &postgres.Config{})
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("bad driver %s", drv)
		}

		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		mig, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s/%s", pwd, p.migrationsDir()), drv, ins)
		if err != nil {
			return err
		}
		mig.Log = &migrateLogger{}
		return f(mig)
	})
}

func (p *Plugin) createDatabase(_ *cli.Context) error {
	drv := viper.GetString("database.driver")
	switch drv {
	case postgresqlDriverName:
		args := viper.GetStringMapString("database.args")
		return web.Shell("psql",
			"-h", args["host"],
			"-p", args["port"],
			"-U", "postgres",
			"-c", fmt.Sprintf(
				"CREATE DATABASE %s WITH ENCODING='UTF8'",
				args["dbname"],
			),
		)
	default:
		return fmt.Errorf("bad driver %s", drv)
	}
}
func (p *Plugin) dropDatabase(_ *cli.Context) error {
	drv := viper.GetString("database.driver")
	switch drv {
	case postgresqlDriverName:
		args := viper.GetStringMapString("database.args")
		return web.Shell("psql",
			"-h", args["host"],
			"-p", args["port"],
			"-U", "postgres",
			"-c", fmt.Sprintf("DROP DATABASE %s", args["dbname"]),
		)
	default:
		return fmt.Errorf("bad driver %s", drv)
	}

}
func (p *Plugin) connectDatabase(_ *cli.Context) error {
	drv := viper.GetString("database.driver")
	switch drv {
	case postgresqlDriverName:
		args := viper.GetStringMapString("database.args")
		return web.Shell("psql",
			"-h", args["host"],
			"-p", args["port"],
			"-U", args["user"],
			args["dbname"],
		)
	default:
		return fmt.Errorf("bad driver %s", drv)
	}
}

func (p *Plugin) startServer(port int, grace bool) error {
	secure := viper.GetBool("server.ssl")
	secret, err := web.SECRET()
	if err != nil {
		return err
	}
	log.Infof(
		"application starting on http://localhost:%d",
		port,
	)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: csrf.Protect(
			secret,
			csrf.Path("/"),
			csrf.Secure(secure),
			csrf.CookieName("csrf"),
			csrf.RequestHeader("Authenticity-Token"),
			csrf.FieldName("authenticity_token"),
		)(p.Router),
	}
	if !grace {
		return srv.ListenAndServe()
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Warn("shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	log.Warn("server exiting")
	return nil
}
