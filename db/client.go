package db

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

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

	tableschemacreation := createschemastatement(outputdb.Tmptargetschema)
	log.Printf("CREATING SCHEMA %s", outputdb.Tmptargetschema)
	if _, err := tx.Exec(tableschemacreation); err != nil {
		log.Fatal("Unable to create db schema ", err)
	}

	schematables := getschemametadata(inputdb, inputdb.Targetschema)

	for _, table := range schematables {

		tablemetadata := gettablemetadata(inputdb, inputdb.Targetschema, table)
		statement := createtablestatement(tablemetadata, table, outputdb.Tmptargetschema)
		log.Printf("CREATING TABLE %s", table)
		if _, err := tx.Exec(statement); err != nil {
			log.Printf("Unable to create table %s  %s", table, err)
			tx.Rollback()

		}

	}

	tx.Commit()
}

func getschemametadata(db Client, tableschema string) []string {
	rows, err := db.Connection.Query(`SELECT table_name 
	FROM information_schema.tables 
	WHERE table_schema=$1`, tableschema)

	if err != nil {
		log.Fatal("Unable to retrieve data from ", err)
	}
	defer rows.Close()
	var schematables []string

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
		}
		schematables = append(schematables, table)
	}

	return schematables
}

func gettablemetadata(db Client, tableschema string, tablename string) [][]string {

	rows, err := db.Connection.Query(`SELECT column_name, data_type 
	                                  FROM information_schema.columns 
									  WHERE table_name =$1 AND table_schema=$2`,
		tablename, tableschema)

	if err != nil {
		log.Fatal("Unable to retrieve schema data ", rows)
	}

	defer rows.Close()
	results := [][]string{}
	for rows.Next() {
		var columnname string
		var datatype string
		var rowresult []string
		if err := rows.Scan(&columnname, &datatype); err != nil {
			log.Fatal(err)
		}
		rowresult = append(rowresult, columnname, datatype)
		results = append(results, rowresult)

	}

	return results
}

func temptableschemaname(tableschema string) string {

	currenttime := strconv.FormatInt(time.Now().Unix(), 10)

	return fmt.Sprintf("%s_temp_%s", tableschema, currenttime)

}

func createtablestatement(tablemetadata [][]string, tablename string, tableschema string) string {

	var fieldstring string
	for _, row := range tablemetadata {
		column := row[0]
		datatype := row[1]
		if datatype == "ARRAY" {
			datatype = "VARCHAR"
		}

		fieldstring += fmt.Sprintf("\"%s\" %s ,", column, datatype)

	}

	return fmt.Sprintf("CREATE TABLE %s.%s (%s)", tableschema, tablename, fieldstring[:len(fieldstring)-1])

}

func createschemastatement(schemaname string) string {

	return fmt.Sprintf("CREATE SCHEMA %s", schemaname)
}

// func copytables(inputdb Client, outputdb Client) {

// 	tx, err := outputdb.Connection.Begin()

// 	if err != nil {
// 		log.Fatal("Unabled to connect to outputdb ", err)
// 	}

// 	rows, err := inputdb.Connection.Query(`SELECT table_name
// 										   FROM information_schema.tables
// 										   WHERE table_schema=$1`, inputdb.Targetschema)

// 	if err != nil {
// 		log.Fatal("Unable to retrieve schema data ", rows)
// 	}

// 	for rows.Next() {
// 		go func ()  {
// 			tx, err := outputdb.Connection.Begin()

// 		}

// 	}
// }
