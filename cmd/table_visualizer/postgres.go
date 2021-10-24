package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

type Postgres struct {
	DB        *sql.DB
	SqlErr    error
	Connected bool
	User      string
	Host      string
	Port      string
	Database  string
}

func (p *Postgres) Connect(user string, password string, host string, database string) {
	hostSplit := strings.Split(host, ":")

	if len(hostSplit) < 2 {
		log.Fatalf("host:port string not valid")
		return
	}

	host = hostSplit[0]
	port := hostSplit[1]
	p.User = user
	p.Host = host
	p.Port = port
	p.Database = database
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, database)
	p.DB, p.SqlErr = sql.Open("postgres", connStr)
	if p.SqlErr != nil {
		p.Connected = false
		log.Fatalf("cannot connect to p database %v", p.SqlErr)
	}
	p.Connected = true
}

func (p *Postgres) Close() {
	p.DB.Close()
}

// Select runs the select query with given args, returns Table struct with cols, rows
func (p *Postgres) Select(query string, args ...interface{}) (Table, error) {
	table := Table{}
	row, err := p.DB.Query(query, args...)
	p.SqlErr = err
	if err != nil {
		return table, err
	}
	defer row.Close()

	columns, err := row.Columns()
	p.SqlErr = err
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
	return table, p.SqlErr
}

func (p *Postgres) GetTable(tableName string) (Table, error) {
	q := "SELECT column_name, data_type, is_nullable, column_default FROM INFORMATION_SCHEMA.COLUMNS where table_name = $1;"
	t, err := p.Select(q, tableName)
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
		fks, _ := p.GetFks(tableName, c[0])
		table.Columns = append(table.Columns, Column{
			Name:     c[0],
			Type:     c[1],
			Nullable: c[2] == "YES",
			// Default:  c[3],
			Fks: fks,
		})
	}

	return table, nil
}

func (p *Postgres) GetTables() ([]Table, error) {
	tables := []Table{}
	q := `SELECT table_name
	FROM information_schema.tables
   	WHERE table_schema='public'
	AND table_type='BASE TABLE';`
	out, err := p.Select(q)
	if err != nil {
		return tables, err
	}

	for _, table := range out.Data {
		t, err := p.GetTable(table[0])
		if err != nil {
			return tables, err
		}
		tables = append(tables, t)
	}

	return tables, nil
}

func (p *Postgres) GetFks(tableName string, colName string) ([]Fk, error) {
	q := `SELECT
		string_agg(ccu.table_name, ',') AS foreign_table_name,
		string_agg(ccu.column_name, ',') AS foreign_column_name
	FROM 
		information_schema.table_constraints AS tc 
	JOIN information_schema.key_column_usage AS kcu
  		ON tc.constraint_name = kcu.constraint_name
  		AND tc.table_schema = kcu.table_schema
		JOIN information_schema.constraint_column_usage AS ccu
  		ON ccu.constraint_name = tc.constraint_name
  		AND ccu.table_schema = tc.table_schema
	WHERE tc.constraint_type = 'FOREIGN KEY' AND tc.table_name = $1 AND kcu.column_name = $2
	GROUP BY tc.table_name;`
	fks := []Fk{}

	out, err := p.Select(q, tableName, colName)
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
