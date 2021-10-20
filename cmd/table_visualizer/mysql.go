package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type MySql struct {
	DB        *sql.DB
	SqlErr    error
	Connected bool
	User      string
	Host      string
	Database  string
}

func (mysql *MySql) Connect(user string, password string, host string, database string) {
	mysql.User = user
	mysql.Host = host
	mysql.Database = database
	connStr := user + ":" + password + "@" + "tcp(" + host + ")/" + database
	mysql.DB, mysql.SqlErr = sql.Open("mysql", connStr)
	if mysql.SqlErr != nil {
		mysql.Connected = false
		log.Fatalf("cannot connect to mysql database %v", mysql.SqlErr)
	}
}
