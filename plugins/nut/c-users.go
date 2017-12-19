package nut

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/SermoDigital/jose/jws"
	"github.com/gin-gonic/gin"
	"github.com/kapmahc/h2o/web"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	gomail "gopkg.in/gomail.v2"
)

func (p *Plugin) deleteUsersSignOut(l string, c *gin.Context) (interface{}, error) {
	user := c.MustGet(CurrentUser).(*User)
	if err := p.Dao.AddLog(p.DB, user.ID, c.ClientIP(), l, "nut.logs.user.sign-out"); err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

func (p *Plugin) getUsersProfile(l string, c *gin.Context) (interface{}, error) {
	user := c.MustGet(CurrentUser).(*User)
	return gin.H{
		"name":  user.Name,
		"email": user.Email,
		"logo":  user.Logo,
	}, nil
}

type fmUserProfile struct {
	Name string `json:"name" binding:"required"`
	Logo string `json:"logo" binding:"required"`
}

func (p *Plugin) postUsersProfile(l string, c *gin.Context) (interface{}, error) {
	var fm fmUserProfile
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	user := c.MustGet(CurrentUser).(*User)
	if err := p.DB.Model(user).Updates(map[string]interface{}{
		"name": fm.Name,
		"logo": fm.Logo,
	}).Error; err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

type fmUserChangePassword struct {
	CurrentPassword      string `json:"currentPassword" binding:"required"`
	NewPassword          string `json:"newPassword" binding:"required,min=6"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"eqfield=NewPassword"`
}

func (p *Plugin) postUsersChangePassword(l string, c *gin.Context) (interface{}, error) {
	var fm fmUserChangePassword
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	user := c.MustGet(CurrentUser).(*User)
	pwd, err := p.Security.Hash([]byte(fm.NewPassword))
	if err != nil {
		return nil, err
	}
	db := p.DB.Begin()
	if err := db.Model(user).Update("password", pwd).Error; err != nil {
		db.Rollback()
		return nil, err
	}
	if err := p.Dao.AddLog(db, user.ID, c.ClientIP(), l, "nut.logs.user.change-password"); err != nil {
		db.Rollback()
		return nil, err
	}
	db.Commit()
	return gin.H{}, nil
}

func (p *Plugin) getUsersLogs(l string, c *gin.Context) (interface{}, error) {
	user := c.MustGet(CurrentUser).(*User)
	var items []Log
	if err := p.DB.Where("user_id = ?", user.ID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

type fmUserSignIn struct {
	Email    string `json:"email" binding:"email"`
	Password string `json:"password" binding:"required"`
}

func (p *Plugin) postUsersSignIn(l string, c *gin.Context) (interface{}, error) {
	var fm fmUserSignIn
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	user, err := p.Dao.GetUserByEmail(p.DB, fm.Email)
	if err != nil {
		return nil, err
	}
	if !p.Security.Check(user.Password, []byte(fm.Password)) {
		return nil, p.I18n.E(l, "nut.errors.user.email-password-not-match")
	}
	if !user.IsConfirm() {
		p.I18n.E(l, "nut.errors.user.not-confirm")
	}
	if user.IsLock() {
		return nil, p.I18n.E(l, "nut.errors.user.is-lock")
	}
	cm := make(jws.Claims)
	cm.Set(UID, user.UID)
	tkn, err := p.Jwt.Sum(cm, time.Hour*24)
	if err != nil {
		return nil, err
	}
	return gin.H{"token": string(tkn)}, nil
}

type fmUserSignUp struct {
	Name                 string `json:"name" binding:"required"`
	Email                string `json:"email" binding:"email"`
	Password             string `json:"password" binding:"required,min=6"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"eqfield=Password"`
}

func (p *Plugin) postUsersSignUp(l string, c *gin.Context) (interface{}, error) {
	log.Printf("%+v", c.Request.Header)
	var fm fmUserSignUp
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	if _, err := p.Dao.GetUserByEmail(p.DB, fm.Email); err == nil {
		return nil, p.I18n.E(l, "nut.errors.user.email-already-exist")
	}
	ip := c.ClientIP()
	tx := p.DB.Begin()
	user, err := p.Dao.AddEmailUser(tx, l, ip, fm.Name, fm.Email, fm.Password)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	if err = p.sendEmail(c, l, user, actConfirm); err != nil {
		log.Error(err)
	}

	return nil, nil
}

type fmUserEmail struct {
	Email string `json:"email" binding:"email"`
}

func (p *Plugin) postUsersConfirm(l string, c *gin.Context) (interface{}, error) {
	var fm fmUserEmail
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}

	user, err := p.Dao.GetUserByEmail(p.DB, fm.Email)
	if err != nil {
		return nil, err
	}
	if user.IsConfirm() {
		return nil, p.I18n.E(l, "nut.errors.user.already-confirm")
	}
	if err := p.sendEmail(c, l, user, actConfirm); err != nil {
		log.Error(err)
	}

	return gin.H{}, nil
}

func (p *Plugin) postUsersUnlock(l string, c *gin.Context) (interface{}, error) {
	var fm fmUserEmail
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}

	user, err := p.Dao.GetUserByEmail(p.DB, fm.Email)
	if err != nil {
		return nil, err
	}
	if !user.IsLock() {
		return nil, p.I18n.E(l, "nut.errors.user.not-lock")
	}
	if err := p.sendEmail(c, l, user, actUnlock); err != nil {
		log.Error(err)
	}

	return gin.H{}, nil
}

func (p *Plugin) postUsersForgotPassword(l string, c *gin.Context) (interface{}, error) {
	var fm fmUserEmail
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}

	user, err := p.Dao.GetUserByEmail(p.DB, fm.Email)
	if err != nil {
		return nil, err
	}
	if err := p.sendEmail(c, l, user, actResetPassword); err != nil {
		log.Error(err)
	}

	return gin.H{}, nil
}

type fmUserResetPassword struct {
	Token                string `json:"token" binding:"required"`
	Password             string `json:"password" binding:"required,min=6"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"eqfield=Password"`
}

func (p *Plugin) postUsersResetPassword(l string, c *gin.Context) (interface{}, error) {
	var fm fmUserResetPassword
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	cm, err := p.Jwt.Validate([]byte(fm.Token))
	if err != nil {
		return nil, err
	}
	if cm.Get("act").(string) != actResetPassword {
		return nil, p.I18n.E(l, "errors.bad-action")
	}
	user, err := p.Dao.GetUserByUID(p.DB, cm.Get("uid").(string))
	if err != nil {
		return nil, err
	}
	pwd, err := p.Security.Hash([]byte(fm.Password))
	if err != nil {
		return nil, err
	}

	tx := p.DB.Begin()
	if err = tx.Model(user).Update("password", pwd).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err = p.Dao.AddLog(tx, user.ID, c.ClientIP(), l, "nut.logs.user.reset-password"); err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return gin.H{}, nil
}

func (p *Plugin) getUsersConfirmToken(l string, c *gin.Context) error {
	cm, err := p.Jwt.Validate([]byte(c.Param("token")))
	if err != nil {
		return err
	}
	if cm.Get("act").(string) != actConfirm {
		return p.I18n.E(l, "errors.bad-action")
	}
	user, err := p.Dao.GetUserByUID(p.DB, cm.Get("uid").(string))
	if err != nil {
		return err
	}
	if user.IsConfirm() {
		return p.I18n.E(l, "nut.errors.user.already-confirm")
	}

	tx := p.DB.Begin()
	if err = tx.Model(user).Update("confirmed_at", time.Now()).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err = p.Dao.AddLog(tx, user.ID, c.ClientIP(), l, "nut.logs.user.confirm"); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	s := p.Layout.Session(c)
	s.AddFlash(p.I18n.T(l, "nut.emails.user.confirm.success"), NOTICE)
	p.Layout.Save(c, s)

	return nil
}

func (p *Plugin) getUsersUnlockToken(l string, c *gin.Context) error {
	cm, err := p.Jwt.Validate([]byte(c.Param("token")))
	if err != nil {
		return err
	}
	if cm.Get("act").(string) != actUnlock {
		return p.I18n.E(l, "errors.bad-action")
	}
	user, err := p.Dao.GetUserByUID(p.DB, cm.Get("uid").(string))
	if err != nil {
		return err
	}
	if !user.IsLock() {
		return p.I18n.E(l, "nut.errors.user.not-lock")
	}

	tx := p.DB.Begin()
	if err = tx.Model(user).Update("locked_at", nil).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err = p.Dao.AddLog(tx, user.ID, c.ClientIP(), l, "nut.logs.unlock"); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	s := p.Layout.Session(c)
	s.AddFlash(p.I18n.T(l, "nut.emails.user.unlock.success"), NOTICE)
	p.Layout.Save(c, s)

	return nil
}

const (
	actConfirm       = "confirm"
	actUnlock        = "unlock"
	actResetPassword = "reset-password"

	// SendEmailJob send email
	SendEmailJob = "send.email"
)

func (p *Plugin) sendEmail(c *gin.Context, lang string, user *User, act string) error {
	cm := jws.Claims{}
	cm.Set("act", act)
	cm.Set("uid", user.UID)
	tkn, err := p.Jwt.Sum(cm, time.Hour*6)
	if err != nil {
		return err
	}

	obj := gin.H{
		"backend":  p.Layout.Backend(c),
		"frontend": p.Layout.Frontend(c),
		"token":    string(tkn),
	}

	subject, err := p.I18n.H(lang, fmt.Sprintf("nut.emails.user.%s.subject", act), obj)
	if err != nil {
		return err
	}
	body, err := p.I18n.H(lang, fmt.Sprintf("nut.emails.user.%s.body", act), obj)
	if err != nil {
		return err
	}

	return p.Jobber.Send(SendEmailJob, 0, map[string]string{
		"to":      user.Email,
		"subject": subject,
		"body":    body,
	})

}

func (p *Plugin) doSendEmail(id string, payload []byte) error {
	var buf bytes.Buffer
	dec := gob.NewDecoder(&buf)
	buf.Write(payload)
	arg := make(map[string]string)
	if err := dec.Decode(&arg); err != nil {
		return err
	}

	to := arg["to"]
	subject := arg["subject"]
	body := arg["body"]
	if viper.GetString("env") != web.PRODUCTION {
		log.Debugf("send to %s: %s\n%s", to, subject, body)
		return nil
	}

	smtp := make(map[string]interface{})
	if err := p.Settings.Get(p.DB, "site.smtp", &smtp); err != nil {
		return err
	}

	sender := smtp["username"].(string)
	msg := gomail.NewMessage()
	msg.SetHeader("From", sender)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	dia := gomail.NewDialer(
		smtp["host"].(string),
		smtp["port"].(int),
		sender,
		smtp["password"].(string),
	)

	return dia.DialAndSend(msg)

}
