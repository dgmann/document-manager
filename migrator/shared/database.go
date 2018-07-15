package shared

import (
	"github.com/jmoiron/sqlx"
	"net/url"
)

type Manager struct {
	Db  *sqlx.DB
	dsn string
}

func NewManager(config Config) *Manager {
	query := url.Values{}
	query.Add("database", config.DbName)
	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(config.Username, config.Password),
		Host:     config.Hostname,
		Path:     config.Instance,
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
