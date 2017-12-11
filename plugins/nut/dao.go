package nut

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kapmahc/h2o/web"
)

// Dao dao
type Dao struct {
	Security *web.Security `inject:""`
	I18n     *web.I18n     `inject:""`
}

// SignIn set sign-in info
func (p *Dao) SignIn(db *gorm.DB, lang, ip, email, password string) (*User, error) {
	user, err := p.GetUserByEmail(db, email)
	if err != nil {
		return nil, err
	}

	if !p.Security.Check([]byte(password), []byte(user.Password)) {
		p.AddLog(db, user.ID, ip, lang, "nut.logs.user.sign-in.failed")
		return nil, p.I18n.E(lang, "nut.errors.user.email-password-not-match")
	}

	if !user.IsConfirm() {
		return nil, p.I18n.E(lang, "nut.errors.user.not-confirm")
	}

	if user.IsLock() {
		return nil, p.I18n.E(lang, "nut.errors.user.is-lock")
	}

	p.AddLog(db, user.ID, ip, lang, "nut.logs.user.sign-in.success")
	user.SignInCount++
	user.LastSignInAt = user.CurrentSignInAt
	user.LastSignInIP = user.CurrentSignInIP
	now := time.Now()
	user.CurrentSignInAt = &now
	user.CurrentSignInIP = ip

	if err = db.Model(user).
		Updates(map[string]interface{}{
			"last_sign_in_at":    user.LastSignInAt,
			"last_sign_in_ip":    user.LastSignInIP,
			"current_sign_in_at": user.CurrentSignInAt,
			"current_sign_in_ip": user.CurrentSignInIP,
			"sign_in_count":      user.SignInCount,
		}).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByUID get user by uid
func (p *Dao) GetUserByUID(db *gorm.DB, uid string) (*User, error) {
	var u User
	if err := db.Where("uid = ?", uid).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// GetUserByEmail get user by email
func (p *Dao) GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	var user User
	if err := db.Where("provider_type = ? AND provider_id = ?", UserTypeEmail, email).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// AddLog add log
func (p *Dao) AddLog(db *gorm.DB, user uint, ip, lang, format string, args ...interface{}) error {
	err := db.Create(&Log{
		UserID:  user,
		IP:      ip,
		Message: p.I18n.T(lang, format, args...),
	}).Error
	return err
}

// AddEmailUser add email user
func (p *Dao) AddEmailUser(db *gorm.DB, lang, ip, name, email, password string) (*User, error) {
	passwd, err := p.Security.Hash([]byte(password))
	if err != nil {
		return nil, err
	}
	user := User{
		Email:           email,
		Password:        passwd,
		Name:            name,
		ProviderType:    UserTypeEmail,
		ProviderID:      email,
		LastSignInIP:    "0.0.0.0",
		CurrentSignInIP: "0.0.0.0",
	}
	user.SetUID()
	user.SetGravatarLogo()

	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}
	p.AddLog(db, user.ID, ip, lang, "nut.logs.user.sign-up")
	return &user, nil
}

//Is is role ?
func (p *Dao) Is(db *gorm.DB, user uint, names ...string) bool {
	for _, name := range names {
		if p.Can(db, user, name, DefaultResourceType, DefaultResourceID) {
			return true
		}
	}
	return false
}

//Can can?
func (p *Dao) Can(db *gorm.DB, user uint, name string, rty string, rid uint) bool {
	var r Role

	if err := db.Where("name = ? AND resource_type = ? AND resource_id = ?", name, rty, rid).
		First(&r).Error; err != nil {
		return false
	}
	var pm Policy
	if err := db.Where("user_id = ? AND role_id = ?").First(&pm).Error; err != nil {
		return false
	}

	return pm.Enable()
}

// GetRole create role if not exist
func (p *Dao) GetRole(db *gorm.DB, name string, rty string, rid uint) (*Role, error) {
	r := Role{}
	err := db.Where("name = ? AND resource_type = ? AND resource_id = ?", name, rty, rid).First(&r).Error
	if err == nil {
		return &r, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	r.Name = name
	r.ResourceID = rid
	r.ResourceType = rty
	r.UpdatedAt = time.Now()
	if err = db.Create(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

//Deny deny permission
func (p *Dao) Deny(db *gorm.DB, role uint, user uint) error {
	err := db.Where("role_id = ? AND user_id = ?", role, user).Delete(Policy{}).Error
	return err
}

// // Authority get roles
// func (p *Dao) Authority(db *gorm.DB, user uint, rty string, rid uint) ([]string, error) {
// 	var items []*Role
//
// 	if err := db.Where("resource_type = ? AND resource_id = ?", rty, rid).
// 		Find(&items).Error; err != nil {
// 		return nil, err
// 	}
// 	var roles []string
// 	for _, r := range items {
// 		var pm Policy
// 		if err := db.Where("role_id = ? AND user_id = ?", r.ID, user).
// 			First(&pm).Error; err == nil {
// 			if pm.Enable() {
// 				roles = append(roles, r.Name)
// 			}
// 		}
// 	}
// 	return roles, nil
// }

// //Allow allow permission
// func (p *Dao) Allow(db *gorm.DB, user uint, role uint, years, months, days int) error {
// 	begin := time.Now()
// 	end := begin.AddDate(years, months, days)
//
// 	var pm Policy
// 	err := db.Where("role_id = ? AND user_id = ?", role, user).First(&pm).Error
// 	if err == nil {
// 		err = db.Model(&pm).
// 			Where("id = ?", pm.ID).Updates(map[string]interface{}{
// 			"start_up":  begin,
// 			"shut_down": end,
// 		}).Error
// 	} else if err == gorm.ErrRecordNotFound {
// 		pm.UserID = user
// 		pm.RoleID = role
// 		pm.StartUp = begin
// 		pm.ShutDown = end
// 		err = db.Create(&pm).Error
// 	}
// 	return err
// }
//
// // ListUserByResource list users by resource
// func (p *Dao) ListUserByResource(role, rty string, rid uint) ([]uint, error) {
// 	ror, err := p.GetRole(o, role, rty, rid)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var ids []uint
// 	var policies []Policy
// 	if _, err := o.QueryTable(new(Policy)).
// 		Filter("role_id", ror.ID).All(&policies); err != nil {
// 		return nil, err
// 	}
// 	for _, pm := range policies {
// 		if pm.Enable() {
// 			ids = append(ids, pm.User.ID)
// 		}
// 	}
// 	return ids, nil
// }
//
// // ListResourcesIds list resource ids by user and role
// func (p *Dao) ListResourcesIds(user uint, role, rty string) ([]uint, error) {
// 	var ids []uint
// 	var policies []Policy
//
// 	if _, err := o.QueryTable(new(Policy)).
// 		Filter("user", user).
// 		All(&policies); err != nil {
// 		return nil, err
// 	}
// 	for _, pm := range policies {
// 		if pm.Enable() {
// 			var ror Role
// 			if err := o.QueryTable(&ror).
// 				Filter("id", pm.Role.ID).
// 				One(&ror); err != nil {
// 				return nil, err
// 			}
// 			if ror.Name == role && ror.ResourceType == rty {
// 				ids = append(ids, ror.ResourceID)
// 			}
// 		}
// 	}
// 	return ids, nil
// }

func (p *Dao) confirmUser(db *gorm.DB, lang, ip string, user uint) error {

	if err := db.Model(&User{}).Where("id = ?", user).Update(
		"confirmed_at", time.Now(),
	).Error; err != nil {
		return err
	}
	return p.AddLog(db, user, ip, lang, "nut.logs.user.confirm")
}

func (p *Dao) setUserPassword(db *gorm.DB, lang, ip string, user uint, password string) error {
	passwd, err := p.Security.Hash([]byte(password))
	if err != nil {
		return err
	}
	if err := db.Model(&User{}).Where("id = ?", user).Update(
		"password", passwd,
	).Error; err != nil {
		return err
	}
	return p.AddLog(db, user, ip, lang, "nut.logs.user.change-password")
}
