package web

// http://www.gnu.org/software/gettext/manual/gettext.html#Language-Codes
import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

const (
	// LOCALE locale key
	LOCALE = "lcoale"
)

// NewI18n create i18n
func NewI18n(path string) (*I18n, error) {
	it := I18n{
		items: make(map[string]string),
	}
	if err := it.loadFromFileSystem(path); err != nil {
		return nil, err
	}
	return &it, nil
}

// Locale locale
type Locale struct {
	tableName struct{}  `sql:"locales"`
	ID        uint      `json:"id"`
	Lang      string    `json:"lang"`
	Code      string    `json:"code"`
	Message   string    `json:"message"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// I18n i18n
type I18n struct {
	db    *gorm.DB
	items map[string]string
}

// Middleware locale middleware
func (p *I18n) Middleware(langs ...string) (gin.HandlerFunc, error) {
	var tags []language.Tag
	for _, l := range langs {
		t, e := language.Parse(l)
		if e != nil {
			return nil, e
		}
		tags = append(tags, t)
	}
	matcher := language.NewMatcher(tags)
	return func(c *gin.Context) {
		lang, written := p.detectLocale(c.Request, LOCALE)
		tag, _, _ := matcher.Match(language.Make(lang))
		if lang != tag.String() {
			written = true
			lang = tag.String()
		}
		if written {
			c.SetCookie(LOCALE, lang, math.MaxInt32, "/", "", c.Request.TLS != nil, false)
		}
	}, nil
}

func (p *I18n) detectLocale(r *http.Request, k string) (string, bool) {
	// 1. Check URL arguments.
	if lang := r.URL.Query().Get(k); lang != "" {
		return lang, true
	}

	// 2. Get language information from cookies.
	if ck, er := r.Cookie(k); er == nil {
		return ck.Value, false
	}

	// 3. Get language information from 'Accept-Language'.
	return r.Header.Get("Accept-Language"), true
}

// Languages language tags
func (p *I18n) Languages(db *gorm.DB) ([]string, error) {
	var langs []string
	err := db.Model(&Locale{}).Select("DISTINCT lang").Scan(&langs).Error
	return langs, err
}

func (p *I18n) loadFromFileSystem(dir string) error {
	const ext = ".ini"
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		name := info.Name()
		if info.IsDir() || filepath.Ext(name) != ext {
			return err
		}
		tag, err := language.Parse(name[:len(name)-len(ext)])
		if err != nil {
			return err
		}
		log.Info("find locale ", tag)
		lang := tag.String()

		cfg, err := ini.Load(path)
		if err != nil {
			return err
		}

		for _, sec := range cfg.Sections() {
			z := sec.Name()
			for k, v := range sec.KeysHash() {
				p.items[lang+"."+z+"."+k] = v
			}
		}

		return nil
	})
}

// Set set
func (p *I18n) Set(db *gorm.DB, lang, code, message string) error {
	var it Locale
	now := time.Now()
	err := db.Select([]string{"id"}).
		Where("lang = ? AND code = ?", lang, code).First(&it).Error
	if err == nil {
		err = db.Model(&it).Update("message", message).Error
	} else if err == gorm.ErrRecordNotFound {
		err = db.Create(&Locale{
			Lang:      lang,
			Code:      code,
			Message:   message,
			UpdatedAt: now,
		}).Error
	}

	if err == nil {
		p.items[lang+"."+code] = message
	}
	return err
}

// H html
func (p *I18n) H(lang, code string, obj interface{}) (string, error) {
	msg, err := p.get(lang, code)
	if err != nil {
		return "", err
	}
	tpl, err := template.New("").Parse(msg)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, obj)
	return buf.String(), err
}

//E error
func (p *I18n) E(lang, code string, args ...interface{}) error {
	msg, err := p.get(lang, code)
	if err != nil {
		return err
	}
	return fmt.Errorf(msg, args...)
}

//T text
func (p *I18n) T(lang, code string, args ...interface{}) string {
	msg, err := p.get(lang, code)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf(msg, args...)
}

func (p *I18n) get(lang, code string) (string, error) {
	var it Locale
	if err := p.db.Select([]string{"message"}).
		Where("lang = ? AND code = ?", lang, code).
		First(&it).Error; err == nil {
		return it.Message, nil
	}
	key := lang + "." + code
	if msg, ok := p.items[key]; ok {
		return msg, nil
	}
	return "", errors.New(key)
}
