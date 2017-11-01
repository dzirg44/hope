package main

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

var DB *sql.DB
var DBS *sql.DB



func MysqlDB(dsn string) (*sql.DB, error) {
db, err := sql.Open("mysql", dsn+"?parseTime=true")
if err != nil {
	return nil, err
}
return db, db.Ping()
}