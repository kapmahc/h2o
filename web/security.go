package web

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"golang.org/x/crypto/bcrypt"
)

// NewSecurity new security
func NewSecurity(key []byte) (*Security, error) {
	cip, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return &Security{cip: cip}, nil
}

// Security security helper
type Security struct {
	cip cipher.Block
}

// Hash ont-way encrypt
func (p *Security) Hash(plain []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(plain, 16)
}

// Check check hash
func (p *Security) Check(encode, plain []byte) bool {
	return bcrypt.CompareHashAndPassword(encode, plain) == nil
}

// Encrypt encrypt
func (p *Security) Encrypt(buf []byte) ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(p.cip, iv)
	val := make([]byte, len(buf))
	cfb.XORKeyStream(val, buf)

	return append(val, iv...), nil
}

// Decrypt decrypt
func (p *Security) Decrypt(buf []byte) ([]byte, error) {
	bln := len(buf)
	cln := bln - aes.BlockSize
	ct := buf[0:cln]
	iv := buf[cln:bln]

	cfb := cipher.NewCFBDecrypter(p.cip, iv)
	val := make([]byte, cln)
	cfb.XORKeyStream(val, ct)
	return val, nil
}
