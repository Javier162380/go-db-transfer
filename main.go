package main

import (
	"fmt"

	"github.com/Javier162380/go-db-transfer/configuration"
	"github.com/Javier162380/go-db-trasfer/cmd"
)

func main() {
	cmd.Init()
	load := configuration.LoadConfiguration(cmd.CONFIGFILE)
	fmt.Printf("%s", load)

}
