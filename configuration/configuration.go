package configuration

import (
	"io/ioutil"
	"log"
	"os"

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

func LoadConfiguration(configfile string) DBCONFIGURATION {

	var dbconfig DBCONFIGURATION
	yamlfile, err := ioutil.ReadFile(configfile)

	if err != nil {
		log.Fatal("Unable to load config file, %s", configfile)
		os.Exit(1)
	}

	err = yaml.Unmarshal(yamlfile, &dbconfig)

	if err != nil {
		log.Fatal("Unable to Unmarshall config file %s error %s", configfile, err)
		os.Exit(1)

	}

	return dbconfig

}
