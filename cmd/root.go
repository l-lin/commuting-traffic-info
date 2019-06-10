package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/l-lin/commuting-traffic-info/config"
	"github.com/l-lin/commuting-traffic-info/twitter"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const cfgFileName = ".commuting-traffic-info"

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "commuting-traffic-info",
	Short: "Check commuting traffic in Paris",
	Run:   run,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Requires commuting line number")
		}
		if _, err := strconv.Atoi(args[0]); err != nil {
			return errors.New("Commuting line number must be an number")
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	lineNb, _ := strconv.Atoi(args[0])
	resultCh := make(chan *twitter.SearchTweetsResult)
	go twitter.SearchTweets(resultCh, lineNb)
	result := <-resultCh
	log.Println(result)
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.commuting-traffic-info.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".commuting-traffic-info" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(cfgFileName)
		cfgFile = fmt.Sprintf("%s/%s.yaml", home, cfgFileName)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// log.Println("Using config file:", viper.ConfigFileUsed())
		key := config.GetAPIKey()
		secretKey := config.GetAPISecretKey()
		if key == "" || secretKey == "" {
			log.Println("Could not read the 'key' and 'secret-key' properties. Initializing it")
			config.InitTwitterAPIKeys(cfgFile)
		}
	} else {
		config.InitTwitterAPIKeys(cfgFile)
	}
}
