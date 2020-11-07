package mysql

import (
	"github.com/bongnv/gwf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Config defines the config for the MYSQL connection.
type Config struct {
	DriverName                string
	DSN                       string
	Conn                      gorm.ConnPool
	SkipInitializeWithVersion bool
	DefaultStringSize         uint
	DisableDatetimePrecision  bool
	DontSupportRenameIndex    bool
	DontSupportRenameColumn   bool
	GormConfig                gorm.Config
}

// WithMYSQL intializes an MySQL instance and registers it to the Application.
func WithMYSQL(cfg Config) gwf.OptionFn {
	return func(app *gwf.Application) {
		gormCfg := mysql.Config{
			DriverName:                cfg.DriverName,
			DSN:                       cfg.DSN,
			Conn:                      cfg.Conn,
			SkipInitializeWithVersion: cfg.SkipInitializeWithVersion,
			DefaultStringSize:         cfg.DefaultStringSize,
			DisableDatetimePrecision:  cfg.DisableDatetimePrecision,
			DontSupportRenameIndex:    cfg.DontSupportRenameIndex,
			DontSupportRenameColumn:   cfg.DontSupportRenameColumn,
		}
		db, err := gorm.Open(mysql.New(gormCfg), &cfg.GormConfig)
		if err != nil {
			panic(err)
		}

		app.MustRegister("db", db)
	}
}

// WithMYSQLByDSN is a short form of WithMYSQL by using dsn only.
func WithMYSQLByDSN(dsn string) gwf.OptionFn {
	return WithMYSQL(Config{
		DSN: dsn,
	})
}
