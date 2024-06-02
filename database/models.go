package database

import "database/sql"

type Database struct {
	Addr  string
	Db    string
	Table string
}

type Connector interface {
	Connect(db Database) *sqlx.Conn
}
