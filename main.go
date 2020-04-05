package main

import (
	"go-db-transfer/cmd"
	"go-db-transfer/configuration"
	"go-db-transfer/db"
)

func main() {
	cmd.Init()
	dbsettings := configuration.LoadConfiguration(cmd.CONFIGFILE)
	inputdb := db.DBClient{Uri: dbsettings.InputDB.URI, Dbtype: dbsettings.InputDB.DBType, Targetschema: dbsettings.InputDB.TargetSchema}
	outputdb := db.DBClient{Uri: dbsettings.OutputDB.URI, Dbtype: dbsettings.OutputDB.DBType, Targetschema: dbsettings.OutputDB.TargetSchema}

	inputdb.Connect()
	outputdb.Connect()

	db.ReplicateSchema(inputdb, outputdb)
}
