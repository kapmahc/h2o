package web

import (
	"crypto/rand"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"syscall"
)

// FuncName get func's name
func FuncName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

// RandomBytes random bytes
func RandomBytes(l int) ([]byte, error) {
	buf := make([]byte, l)
	if _, err := rand.Read(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

//Shell exec shell command
func Shell(cmd string, args ...string) error {
	bin, err := exec.LookPath(cmd)
	if err != nil {
		return err
	}
	return syscall.Exec(bin, append([]string{cmd}, args...), os.Environ())
}
