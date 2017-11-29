package web

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type loggedResponse struct {
	http.ResponseWriter
	status int
}

func (l *loggedResponse) WriteHeader(status int) {
	l.status = status
	l.ResponseWriter.WriteHeader(status)
}

// LoggerMiddleware logger middleware
func LoggerMiddleware(wrt http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	log.Infof("%s %s %s", req.Proto, req.Method, req.URL)
	begin := time.Now()
	lw := loggedResponse{ResponseWriter: wrt}
	next(&lw, req)
	log.Infof("%d %s", lw.status, time.Now().Sub(begin))
}
