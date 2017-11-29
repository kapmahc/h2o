package nut

import "log"

type migrateLogger struct {
}

// Printf is like fmt.Printf
func (p *migrateLogger) Printf(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Verbose should return true when verbose logging output is wanted
func (p *migrateLogger) Verbose() bool { return true }
