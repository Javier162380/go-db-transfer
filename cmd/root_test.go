package cmd

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

var testconfigfile string

func RootCommand() *cobra.Command {
	var root = &cobra.Command{
		Use:   "go-db-transfer",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
    examples and usage of using your application. For example:
    Cobra is a CLI library for Go that empowers applications.
    This application is a tool to generate the needed files
    to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			configfile, _ := cmd.Flags().GetString("config")
			if strings.HasSuffix(configfile, ".yaml") == false {
				log.Fatal("Invalid file extension only Yaml it is supported")
				os.Exit(1)
			}

		},
	}

	root.Flags().StringVar(&testconfigfile, "config", "", "Config file for the app")

	return root
}

func TestCli(t *testing.T) {

	cmd := RootCommand()
	cmd.SetArgs([]string{"--config", "testconf.yaml"})
	cmd.Execute()

	assert.Equal(t, testconfigfile, "testconf.yaml")

}
