package configuration

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type DBCONFIGURATION struct {
	InputDB  DBSETTINGS `yaml:"input_db"`
	OutputDB DBSETTINGS `yaml:"output_db"`
}

type DBSETTINGS struct {
	URI          string `yaml:"uri"`
	DBType       string `yaml:"db_type"`
	TargetSchema string `yaml:"target_schema"`
}

// LoadConfiguration method to load cli config yaml into a struct
func LoadConfiguration(configfile string) DBCONFIGURATION {

	var dbconfig DBCONFIGURATION
	yamlfile, err := ioutil.ReadFile(configfile)

	if err != nil {
		log.Fatal("Unable to load config file, ", configfile)

	}

	err = yaml.Unmarshal(yamlfile, &dbconfig)

	if err != nil {
		log.Fatal("Unable to Unmarshall config file error ", configfile, err)

	}

	return dbconfig

}
