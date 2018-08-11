package shared

import (
	"github.com/jmoiron/sqlx"
	"net/url"
)

type Manager struct {
	Db  *sqlx.DB
	dsn string
}

func NewManager(dbName, username, password, hostname, instance string) *Manager {
	query := url.Values{}
	query.Add("database", dbName)
	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(username, password),
		Host:     hostname,
		Path:     instance,
		RawQuery: query.Encode(),
	}

	return &Manager{dsn: u.String()}
}

func (m *Manager) Open() error {
	db, err := sqlx.Connect("sqlserver", m.dsn)
	if err != nil {
		return err
	}
	m.Db = db
	return nil
}

func (m *Manager) Close() {
	m.Db.Close()
}
