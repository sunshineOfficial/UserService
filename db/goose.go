package db

import (
	"fmt"
	"io/fs"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

type gooseLogger struct {
	log *zap.Logger
}

func (l gooseLogger) Printf(format string, v ...interface{}) {
	l.log.Debug(fmt.Sprintf(format, v...))
}

func (l gooseLogger) Fatalf(format string, v ...interface{}) {
	l.log.Fatal(fmt.Sprintf(format, v...))
}

func Migrate(fs fs.FS, log *zap.Logger, db *sqlx.DB, path string) error {
	goose.SetVerbose(false)
	goose.SetLogger(gooseLogger{log: log})
	goose.SetBaseFS(fs)
	goose.SetTableName("db_version")

	return goose.Up(db.DB, path)
}
