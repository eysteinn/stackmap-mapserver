/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"sidecar/internal/pkg/fetch"
	"sidecar/internal/pkg/logger"
	"sidecar/internal/pkg/types"
	"time"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sidecar",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		run() //format.Dotest1()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sidecar.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.Flags().StringP("product", "p", "", "Product to base mapfile on.")

}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func run() {
	sqldata := types.SQLData{}

	sqldata.SQLUser = getEnv("PSQLUSER", "postgres")
	sqldata.SQLPass = getEnv("PSQLPASS", "")
	sqldata.SQLDB = getEnv("PSQLDB", "postgres")
	sqldata.SQLHost = getEnv("PSQLHOST", "psqlapi-service.default.svc.cluster.local")
	apiHost := getEnv("APIHOST", "")

	err := fetch.FetchAllProducts("./mapfiles/", apiHost, sqldata)
	if err != nil {
		logger.GetLogger().Panic(err)
	}

	fmt.Println("Going to sleep")
	for {
		time.Sleep(time.Second)
	}
	/*err = rabbitmq.DoRabbitMQ()
	if err != nil {
		logger.GetLogger().Panic(err)
	}*/
}
