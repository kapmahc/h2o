package web

import (
	"net/http"
	"time"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/SermoDigital/jose/jwt"
)

// NewJwt new jwt
func NewJwt(key []byte, method crypto.SigningMethod) *Jwt {
	return &Jwt{key: key, method: method}
}

// Jwt jwt
type Jwt struct {
	key    []byte
	method crypto.SigningMethod
}

//Validate check jwt
func (p *Jwt) Validate(buf []byte) (jwt.Claims, error) {
	tk, err := jws.ParseJWT(buf)
	if err != nil {
		return nil, err
	}
	if err = tk.Validate(p.key, p.method); err != nil {
		return nil, err
	}
	return tk.Claims(), nil
}

// Parse parse
func (p *Jwt) Parse(r *http.Request) (jwt.Claims, error) {
	tk, err := jws.ParseJWTFromRequest(r)
	if err != nil {
		return nil, err
	}
	if err = tk.Validate(p.key, p.method); err != nil {
		return nil, err
	}
	return tk.Claims(), nil
}

//Sum create jwt token
func (p *Jwt) Sum(cm jws.Claims, exp time.Duration) ([]byte, error) {
	now := time.Now()
	cm.SetNotBefore(now)
	cm.SetExpiration(now.Add(exp))

	jt := jws.NewJWT(cm, p.method)
	return jt.Serialize(p.key)
}
