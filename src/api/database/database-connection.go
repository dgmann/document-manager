package database

import (
	"github.com/gocraft/dbr"
	log "github.com/sirupsen/logrus"
)

func NewPostgresConnection(dsn string) *dbr.Connection {
	conn, _ := dbr.Open("postgres", dsn, nil)
	err := conn.Ping()
	if err != nil {
		log.WithFields(log.Fields{
			"dialect": "postgres",
			"dsn":     "dsn",
			"error":   err,
		}).Fatal("Failed to connect to database")
	}
	return conn
}
