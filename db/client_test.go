package db

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

var dbc = Client{URI: "postgres://postgres:postgres@localhost:5434/postgres?sslmode=disable", Dbtype: "postgres",
	Targetschema: "public"}

func TestConnect(t *testing.T) {

	dbc.Connect()

	assert.Contains(t, dbc.Tmptargetschema, "public_temp")
}
func TestTempSchema(t *testing.T) {

	schema := temptableschemaname("public")

	fmt.Printf("schema")
	assert.Contains(t, schema, "public_temp")

}

func TestCreateTableStatement(t *testing.T) {
	tablemetadata := [][]string{[]string{"field1", "text"}, []string{"field2", "integer"}}

	statement := createtablestatement(tablemetadata, "test", "public")

	assert.Equal(t, statement, "CREATE TABLE public.test (\"field1\" text ,\"field2\" integer )")

}

func TestCreateSchemaStatement(t *testing.T) {

	statement := createschemastatement("public")

	assert.Equal(t, statement, "CREATE SCHEMA public")
}

func TestGetTableMetadata(t *testing.T) {

	dbc.Connect()

	tx, err := dbc.Connection.Begin()

	if err != nil {
		t.Fatal("Unable to connect to a testdb, ", err)
	}

	if _, err := tx.Exec("DROP TABLE IF EXISTS public.temp_test;CREATE TABLE public.temp_test (a text, c integer)"); err != nil {
		t.Fatal("Unable to create seed tables ", err)
	}

	tx.Commit()

	results := gettablemetadata(dbc, "public", "temp_test")

	assert.Equal(t, results, [][]string{[]string{"a", "text"}, []string{"c", "integer"}})

}

func TestGetSchemaMetadata(t *testing.T) {

	dbc.Connect()

	tx, err := dbc.Connection.Begin()

	if err != nil {
		t.Fatal("Unable to connect to a testdb, ", err)
	}

	seedstatements := "DROP SCHEMA IF EXISTS test_temp CASCADE; CREATE SCHEMA test_temp; CREATE TABLE test_temp.test (a text)"

	if _, err := tx.Exec(seedstatements); err != nil {
		log.Fatal("Unable to create seed data", err)
	}

	tx.Commit()

	results := getschemametadata(dbc, "test_temp")

	assert.Equal(t, results, []string{"test"})

}
