package mariadb

// import (
// 	"database/sql"
// 	"time"

// 	"go.elastic.co/apm/module/apmsql"
// 	_ "go.elastic.co/apm/module/apmsql/mysql"
// )

// var mariadb *sql.DB
// var driver = "mysql"

// type MariadbConfig struct {
// 	DSN             string
// 	ConnMaxLifetime time.Duration
// 	MaxOpenConns    int
// 	MaxIdleConns    int
// }

// func GetMariadb(config MariadbConfig) (*sql.DB, error) {
// 	if mariadb != nil {
// 		return mariadb, nil
// 	}

// 	db, err := apmsql.Open(driver, config.DSN)
// 	if err != nil {
// 		return db, err
// 	}

// 	db.SetConnMaxLifetime(config.ConnMaxLifetime)
// 	db.SetMaxOpenConns(config.MaxOpenConns)
// 	db.SetMaxIdleConns(config.MaxIdleConns)

// 	mariadb = db

// 	return mariadb, nil
// }
