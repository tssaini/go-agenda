package sqldb

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// DB sql db connection
type DB struct {
	*sql.DB
}

// NewDB create new db connection
func NewDB(dbType, dbConnString string) (*DB, error) {
	db, err := sql.Open(dbType, dbConnString)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
