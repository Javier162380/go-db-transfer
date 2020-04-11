package configuration

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestConfiguration(t *testing.T) {
     
    configuration := LoadConfiguration("../config.yaml")

    expectedconfiguration := DBCONFIGURATION{
        InputDB: DBSETTINGS{URI: "postgres://postgres:postgres@localhost:5434/postgres?sslmode=disable", 
                            DBType: "postgres", 
                            TargetSchema: "information_schema"},
        OutputDB: DBSETTINGS{URI: "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable", 
                                 DBType: "postgres", 
                                 TargetSchema: "public"},
    }

    assert.Equal(t, configuration, expectedconfiguration)
}