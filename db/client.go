package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type DBClient struct {
	Uri          string
	Dbtype       string
	Targetschema string
	Connection   *sql.DB
}

func (dbc *DBClient) Connect() {

	switch {
	case dbc.Dbtype == "postgres":
		var err error
		dbc.Connection, err = sql.Open("postgres", dbc.Uri)

		if err != nil {
			log.Fatal("Unable to connect to DB %s", err)
			os.Exit(1)
		}

	case dbc.Dbtype != "postgres":
		log.Fatal("DB type not supportted currently")

	}
}

func ReplicateSchema(input_db DBClient, output_db DBClient) {

	rows, err := input_db.Connection.Query("SELECT table_name FROM information_schema.tables WHERE table_schema=$1", input_db.Targetschema)

	if err != nil {
		log.Fatal("Unable to retrieve schema data %s", rows)
		os.Exit(1)
	}

	defer rows.Close()
	for rows.Next() {
		var table_name string
		if err := rows.Scan(&table_name); err != nil {
			log.Fatal(err)
		}
		log.Printf("table name %s\n", table_name)

	}
}
