package main

import (
	"flag"
	"log"
	"os"

	"github.com/goccy/go-graphviz"
)

func main() {
	user := flag.String("u", "", "username of the database")
	password := flag.String("p", "", "password of the database user")
	host := flag.String("h", "", "host and port ex. 127.0.0.1:3306")
	database := flag.String("d", "", "database name")
	rdbms := flag.String("s", "", "RDBMS mysql, mariadb, postgres")
	tableName := flag.String("t", "", "table name")
	showAll := flag.Bool("a", true, "show all tables")
	outputPath := flag.String("o", "", "output path /home/user/Desktop/filename.jpg")

	flag.Parse()

	if len(*user) == 0 || len(*password) == 0 || len(*host) == 0 || len(*database) == 0 || len(*rdbms) == 0 || len(*outputPath) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	var db Database

	switch *rdbms {
	case "mysql", "mariadb":
		db = &MySql{}
	case "postgres":
		db = &Postgres{}
	default:
	}

	db.Connect(*user, *password, *host, *database)
	defer db.Close()
	var out string
	tables := []Table{}

	if *showAll && len(*tableName) == 0 {
		tables, _ = db.GetTables()
	} else {
		table, _ := db.GetTable(*tableName)
		tables = append(tables, table)
		tables = collectFkTables(&table, db, tables)
	}

	out = generateGraph(tables)

	g := graphviz.New()
	graph, err := graphviz.ParseBytes([]byte(out))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := graph.Close(); err != nil {
			log.Fatal(err)
		}
		g.Close()
	}()

	if err := g.RenderFilename(graph, graphviz.JPG, *outputPath); err != nil {
		log.Fatal(err)
	}

	if err := g.RenderFilename(graph, graphviz.Format(graphviz.DOT), *outputPath+".dot"); err != nil {
		log.Fatal(err)
	}
}

func collectFkTables(table *Table, db Database, tables []Table) []Table {
	for _, c := range table.Columns {
		for _, fk := range c.Fks {
			if !isExists(fk.Table, tables) {
				fkT, err := db.GetTable(fk.Table)
				if err != nil {
					log.Fatal(err)
				}
				tables = collectFkTables(&fkT, db, tables)
				tables = append(tables, fkT)
			}
		}
	}
	return tables
}

func isExists(tableName string, tables []Table) bool {
	for _, t := range tables {
		if t.Name == tableName {
			return true
		}
	}
	return false
}
