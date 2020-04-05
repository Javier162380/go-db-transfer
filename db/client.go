package db

import (
	"database/sql"
	"fmt"
	"log"

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
	case dbc.Dbtype == "postrges":
		conn, err := sql.Open("postgres", dbc.Uri)

		if err != nil {
			log.Fatal("Unable to connect to DB %s", err)
		}
		dbc.Connection = conn

	case dbc.Dbtype != "postgres":
		log.Fatal("DB type not supportted currently")

	}
}

func ReplicateSchema(input_db DBClient, output_db DBClient) {

	rows, err := input_db.Connection.Query("SELECT table_name FROM information_schema.tables WHERE table_schema=$1", input_db.Targetschema)

	if err != nil {

	}

	fmt.Printf("%s", rows)
}
