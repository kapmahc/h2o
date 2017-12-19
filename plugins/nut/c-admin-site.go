package nut

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/user"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/kapmahc/h2o/web"
	"github.com/spf13/viper"
)

type fmSiteHome struct {
	Favicon string `json:"favicon" binding:"required"`
	Theme   string `json:"theme" binding:"required"`
}

func (p *Plugin) postAdminSiteHome(l string, c *gin.Context) (interface{}, error) {
	var fm fmSiteHome
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	db := p.DB.Begin()
	for k, v := range map[string]string{
		"site.favicon":    fm.Favicon,
		"site.home.theme": fm.Theme,
	} {
		if err := p.Settings.Set(db, k, v, false); err != nil {
			db.Rollback()
			return nil, err
		}
	}
	db.Commit()
	return gin.H{}, nil
}

func (p *Plugin) getAdminSiteSMTP(l string, c *gin.Context) (interface{}, error) {
	var smtp map[string]interface{}
	if err := p.Settings.Get(p.DB, "site.smtp", &smtp); err == nil {
		delete(smtp, "password")
	} else {
		smtp = map[string]interface{}{
			"host": "localhost",
			"port": 25,
			"user": "whoami@change-me.com",
		}
	}
	return smtp, nil
}

type fmSiteSMTP struct {
	Host                 string `json:"host" binding:"required"`
	Port                 int    `json:"port"`
	User                 string `json:"user" binding:"email"`
	Password             string `json:"password" binding:"required,min=6"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"eqfield=Password"`
}

func (p *Plugin) postAdminSiteSMTP(l string, c *gin.Context) (interface{}, error) {
	var fm fmSiteSMTP
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	if err := p.Settings.Set(p.DB, "site.smtp", map[string]interface{}{
		"host":     fm.Host,
		"port":     fm.Port,
		"user":     fm.User,
		"password": fm.Password,
	}, true); err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

func (p *Plugin) getAdminSiteSeo(l string, c *gin.Context) (interface{}, error) {
	var googleVerifyCode string
	p.Settings.Get(p.DB, "site.google.verify.code", &googleVerifyCode)
	var baiduVerifyCode string
	p.Settings.Get(p.DB, "site.baidu.verify.code", &baiduVerifyCode)
	return gin.H{
		"googleVerifyCode": googleVerifyCode,
		"baiduVerifyCode":  baiduVerifyCode,
	}, nil
}

type fmSiteSeo struct {
	GoogleVerifyCode string `json:"googleVerifyCode" binding:"required"`
	BaiduVerifyCode  string `json:"baiduVerifyCode" binding:"required"`
}

func (p *Plugin) postAdminSiteSeo(l string, c *gin.Context) (interface{}, error) {
	var fm fmSiteSeo
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	db := p.DB.Begin()
	for k, v := range map[string]string{
		"site.google.verify.code": fm.GoogleVerifyCode,
		"site.baidu.verify.code":  fm.BaiduVerifyCode,
	} {
		if err := p.Settings.Set(db, k, v, false); err != nil {
			db.Rollback()
			return nil, err
		}
	}
	db.Commit()
	return gin.H{}, nil
}

type fmSiteAuthor struct {
	Email string `json:"email" binding:"email"`
	Name  string `json:"name" binding:"required"`
}

func (p *Plugin) postAdminSiteAuthor(l string, c *gin.Context) (interface{}, error) {
	var fm fmSiteAuthor
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	db := p.DB.Begin()
	for k, v := range map[string]string{
		"email": fm.Email,
		"name":  fm.Name,
	} {
		if err := p.Settings.Set(db, "site.author."+k, v, false); err != nil {
			db.Rollback()
			return nil, err
		}
	}
	db.Commit()
	return gin.H{}, nil
}

type fmSiteInfo struct {
	Title       string `json:"title" binding:"required"`
	Subhead     string `json:"subhead" binding:"required"`
	Keywords    string `json:"keywords" binding:"required"`
	Description string `json:"description" binding:"required"`
	Copyright   string `json:"copyright" binding:"required"`
}

func (p *Plugin) postAdminSiteInfo(l string, c *gin.Context) (interface{}, error) {
	var fm fmSiteInfo
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	db := p.DB.Begin()
	for k, v := range map[string]string{
		"title":       fm.Title,
		"subhead":     fm.Subhead,
		"keywords":    fm.Keywords,
		"description": fm.Description,
		"copyright":   fm.Copyright,
	} {
		if err := p.I18n.Set(db, l, "site."+k, v); err != nil {
			db.Rollback()
			return nil, err
		}
	}
	db.Commit()
	return gin.H{}, nil
}

func (p *Plugin) getAdminSiteStatus(l string, c *gin.Context) (interface{}, error) {
	ret := gin.H{
		"jobber": p.Jobber.Status(),
		"routes": p.Router.Routes(),
	}
	var err error
	if ret["os"], err = p._osStatus(); err != nil {
		return nil, err
	}
	if ret["network"], err = p._networkStatus(); err != nil {
		return nil, err
	}
	if ret["database"], err = p._databaseStatus(); err != nil {
		return nil, err
	}
	if ret["redis"], err = p._redisStatus(); err != nil {
		return nil, err
	}
	return ret, nil
}
func (p *Plugin) _osStatus() (gin.H, error) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	hn, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	hu, err := user.Current()
	if err != nil {
		return nil, err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	var ifo syscall.Sysinfo_t
	if err := syscall.Sysinfo(&ifo); err != nil {
		return nil, err
	}
	return gin.H{
		"app author":           fmt.Sprintf("%s <%s>", web.AuthorName, web.AuthorEmail),
		"app licence":          web.Copyright,
		"app version":          fmt.Sprintf("%s(%s)", web.Version, web.BuildTime),
		"app root":             pwd,
		"who-am-i":             fmt.Sprintf("%s@%s", hu.Username, hn),
		"go version":           runtime.Version(),
		"go root":              runtime.GOROOT(),
		"go runtime":           runtime.NumGoroutine(),
		"go last gc":           time.Unix(0, int64(mem.LastGC)).Format(time.ANSIC),
		"os cpu":               runtime.NumCPU(),
		"os ram(free/total)":   fmt.Sprintf("%dM/%dM", ifo.Freeram/1024/1024, ifo.Totalram/1024/1024),
		"os swap(free/total)":  fmt.Sprintf("%dM/%dM", ifo.Freeswap/1024/1024, ifo.Totalswap/1024/1024),
		"go memory(alloc/sys)": fmt.Sprintf("%dM/%dM", mem.Alloc/1024/1024, mem.Sys/1024/1024),
		"os time":              time.Now().Format(time.ANSIC),
		"os arch":              fmt.Sprintf("%s(%s)", runtime.GOOS, runtime.GOARCH),
		"os uptime":            (time.Duration(ifo.Uptime) * time.Second).String(),
		"os loads":             ifo.Loads,
		"os procs":             ifo.Procs,
	}, nil
}
func (p *Plugin) _networkStatus() (gin.H, error) {
	sts := gin.H{}
	ifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, v := range ifs {
		ips := []string{v.HardwareAddr.String()}
		adrs, err := v.Addrs()
		if err != nil {
			return nil, err
		}
		for _, adr := range adrs {
			ips = append(ips, adr.String())
		}
		sts[v.Name] = ips
	}
	return sts, nil
}

func (p *Plugin) _databaseStatus() (gin.H, error) {
	val := gin.H{
		"drivers": strings.Join(sql.Drivers(), ", "),
	}
	db := p.DB.DB()
	args := viper.GetStringMap("database.args")
	switch viper.GetString("database.driver") {
	case "postgres":
		var version string
		if err := db.QueryRow("select version()").Scan(&version); err != nil {
			return nil, err
		}
		val["version"] = version

		// http://blog.javachen.com/2014/04/07/some-metrics-in-postgresql.html
		var size string
		if err := db.QueryRow("select pg_size_pretty(pg_database_size('postgres'))").Scan(&size); err != nil {
			return nil, err
		}
		val["size"] = size

		rows, err := db.Query("select pid,current_timestamp - least(query_start,xact_start) AS runtime,substr(query,1,25) AS current_query from pg_stat_activity where not pid=pg_backend_pid()")
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var pid int
			var ts time.Time
			var qry string
			rows.Scan(&pid, &ts, &qry)
			val[fmt.Sprintf("pid-%d", pid)] = fmt.Sprintf("%s (%v)", ts.Format("15:04:05.999999"), qry)
		}

		val["url"] = fmt.Sprintf("%s@%s:%d/%s", args["user"], args["host"], args["port"], args["dbname"])
	}
	return val, nil
}

func (p *Plugin) _redisStatus() (string, error) {
	c := p.Redis.Get()
	defer c.Close()
	return redis.String(c.Do("INFO"))
}
