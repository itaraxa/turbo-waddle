package postgres

import "errors"

var (
	ErrOpenConnection   = errors.New("NewPostgresRepository: error openning connection")
	ErrMigration        = errors.New("NewPostgresRepository: error migration database")
	ErrInitTransaction  = errors.New("NewTransaction: error init new transaction")
	ErrAddUserQueryToDB = errors.New("AddUser: query to database error")
)
