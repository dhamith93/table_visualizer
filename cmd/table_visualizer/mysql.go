package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

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
	mysql.Connected = true
}

func (mysql *MySql) Close() {
	mysql.DB.Close()
}

func (mysql *MySql) Select(query string, args ...interface{}) (Table, error) {
	table := Table{}
	row, err := mysql.DB.Query(query, args...)
	mysql.SqlErr = err
	if err != nil {
		return table, err
	}
	defer row.Close()

	columns, err := row.Columns()
	mysql.SqlErr = err
	if err != nil {
		return table, err
	}

	cols := []Column{}

	for _, c := range columns {
		cols = append(cols, Column{Name: c})
	}

	output := make([][]string, 0)
	rawResult := make([][]byte, len(columns))
	dest := make([]interface{}, len(columns))
	for i := range rawResult {
		dest[i] = &rawResult[i]
	}

	for row.Next() {
		row.Scan(dest...)
		res := make([]string, 0)
		for _, raw := range rawResult {
			if raw != nil {
				res = append(res, string(raw))
			}
		}
		output = append(output, res)
	}

	table.Columns = cols
	table.Data = output
	return table, mysql.SqlErr
}

func (mysql *MySql) GetTable(tableName string) (Table, error) {
	// not using parameterized query cause it is not working with DESC query
	q := "DESC " + tableName
	t, err := mysql.Select(q)
	if err != nil {
		return t, err
	}

	if len(t.Data) == 0 {
		return t, fmt.Errorf("table not found")
	}

	table := Table{
		Name:    tableName,
		Columns: []Column{},
	}

	for _, c := range t.Data {
		fks, _ := mysql.GetFks(tableName, c[0])
		table.Columns = append(table.Columns, Column{
			Name:     c[0],
			Type:     c[1],
			Nullable: c[2] == "YES",
			Key:      c[3],
			Default:  c[4],
			Fks:      fks,
		})
	}
	return table, nil
}

func (mysql *MySql) GetTables() ([]Table, error) {
	tables := []Table{}
	out, err := mysql.Select("SHOW TABLES;")
	if err != nil {
		return tables, err
	}

	for _, table := range out.Data {
		t, err := mysql.GetTable(table[0])
		if err != nil {
			return tables, err
		}
		tables = append(tables, t)
	}

	return tables, nil
}

func (mysql *MySql) GetFks(tableName string, colName string) ([]Fk, error) {
	q := `SELECT GROUP_CONCAT(TABLE_NAME), GROUP_CONCAT(COLUMN_NAME), REFERENCED_COLUMN_NAME FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE REFERENCED_TABLE_SCHEMA = ? AND REFERENCED_TABLE_NAME = ? AND REFERENCED_COLUMN_NAME = ? GROUP BY REFERENCED_COLUMN_NAME;`
	fks := []Fk{}

	out, err := mysql.Select(q, mysql.Database, tableName, colName)
	if err != nil {
		return fks, err
	}

	for _, row := range out.Data {
		refTables := strings.Split(row[0], ",")
		refCols := strings.Split(row[1], ",")

		for i, rt := range refTables {
			fks = append(fks, Fk{
				Table:  rt,
				RefCol: refCols[i],
			})
		}
	}

	return fks, nil
}
