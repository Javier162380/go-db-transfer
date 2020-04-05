package db

import (
	"database/sql"
	"fmt"
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

	tmptableschema := temptableschemaname(output_db.Targetschema)

	tx, err := output_db.Connection.Begin()

	if err != nil {
		log.Fatal("Unabled to connect to outputdb", err)
		os.Exit(1)
	}
	tx.Exec("CREATE SCHEMA $1", tmptableschema)

	rows, err := input_db.Connection.Query("SELECT table_name FROM information_schema.tables WHERE table_schema=$1", input_db.Targetschema)

	if err != nil {
		log.Fatal("Unable to retrieve schema data %s", rows)
		os.Exit(1)
	}

	defer rows.Close()
	for rows.Next() {
		var tablename string
		if err := rows.Scan(&tablename); err != nil {
			log.Fatal(err)
		}

		tablemetadata := gettablemetadata(input_db, input_db.Targetschema, tablename)
		statement := createtablestatement(tablemetadata, tablename, tmptableschema)
		tx.Exec(statement)

	}

	tx.Commit()
}

func gettablemetadata(db DBClient, tableschema string, tablename string) map[string]string {

	rows, err := db.Connection.Query("SELECT column_name, data_type FROM information_schema.columns WHERE table_name =$1 AND table_schema=$2", tablename, tableschema)

	if err != nil {
		log.Fatal("Unable to retrieve schema data %s", rows)
		os.Exit(1)
	}

	defer rows.Close()
	results := make(map[string]string)
	for rows.Next() {
		var columnname string
		var datatype string

		if err := rows.Scan(&columnname, &datatype); err != nil {
			log.Fatal(err)
		}

		results[columnname] = datatype

	}

	return results
}

func temptableschemaname(tableschema string) string {

	return fmt.Sprintf("%s_temp", tableschema)

}
func createtablestatement(tablemetadata map[string]string, tablename string, tableschema string) string {

	var fieldstring string
	for column, datetype := range tablemetadata {
		fieldstring += fmt.Sprintf("\"%s\" %s, ", column, datetype)

	}

	return fmt.Sprintf("CREATE TABLE %s.%s (%s)", tableschema, tablename, fieldstring)

}
