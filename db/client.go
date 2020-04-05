package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	//postgresql driver
	_ "github.com/lib/pq"
)

// Client main db fields
type Client struct {
	URI             string
	Dbtype          string
	Targetschema    string
	Connection      *sql.DB
	Tmptargetschema string
}

// Connect to sql db
func (dbc *Client) Connect() {

	switch {
	case dbc.Dbtype == "postgres":
		var err error
		dbc.Connection, err = sql.Open("postgres", dbc.URI)

		if err != nil {
			log.Fatal("Unable to connect to DB ", err)
			os.Exit(1)
		}

		dbc.Tmptargetschema = temptableschemaname(dbc.Targetschema)

	case dbc.Dbtype != "postgres":
		log.Fatal("DB type not supportted currently")

	}
}

// ReplicateSchema a schema between to databases
func ReplicateSchema(inputdb Client, outputdb Client) {

	tx, err := outputdb.Connection.Begin()

	if err != nil {
		log.Fatal("Unabled to connect to outputdb ", err)
	}

	tableschemacreation := fmt.Sprintf("CREATE SCHEMA %s", outputdb.Tmptargetschema)

	if _, err := tx.Exec(tableschemacreation); err != nil {
		log.Fatal("Unable to create db schema ", err)
	}

	rows, err := inputdb.Connection.Query(`SELECT table_name 
										   FROM information_schema.tables 
										   WHERE table_schema=$1`, inputdb.Targetschema)

	if err != nil {
		log.Fatal("Unable to retrieve schema data ", err)
	}

	defer rows.Close()
	for rows.Next() {
		var tablename string
		if err := rows.Scan(&tablename); err != nil {
			log.Fatal(err)
		}

		tablemetadata := gettablemetadata(inputdb, inputdb.Targetschema, tablename)
		statement := createtablestatement(tablemetadata, tablename, outputdb.Tmptargetschema)
		if _, err := tx.Exec(statement); err != nil {
			log.Printf("Unable to create table %s  %s", tablename, err)
			tx.Rollback()
		}

	}

	tx.Commit()
}

func gettablemetadata(db Client, tableschema string, tablename string) map[string]string {

	rows, err := db.Connection.Query(`SELECT column_name, data_type 
	                                  FROM information_schema.columns 
									  WHERE table_name =$1 AND table_schema=$2`,
		tablename, tableschema)

	if err != nil {
		log.Fatal("Unable to retrieve schema data ", rows)
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
		fieldstring += fmt.Sprintf("\"%s\" %s ,", column, datetype)

	}

	return fmt.Sprintf("CREATE TABLE %s.%s (%s)", tableschema, tablename, fieldstring[:len(fieldstring)-1])

}
