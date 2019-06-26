package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/kyokomi/emoji"
	"github.com/l-lin/commuting-traffic-info/config"
	"github.com/l-lin/commuting-traffic-info/format"
	"github.com/l-lin/commuting-traffic-info/traffic"
	"github.com/l-lin/commuting-traffic-info/twitter"
	"github.com/logrusorgru/aurora"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const cfgFileName = ".commuting-traffic-info"

var cfgFile string
var formatType string

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
	filteredTweets := twitter.FilterTweets(result.Tweets)
	s := traffic.GetStatus(filteredTweets)
	var formatter format.Formatter
	if "json" == formatType {
		formatter = &format.JSONFormatter{}
	} else if "rocket.chat" == formatType {
		webhook := getRocketChatWebhook()
		formatter = &format.RocketChatFormatter{Webhook: webhook}
	} else {
		formatter = &format.ConsoleFormatter{}
	}
	formatter.Format(lineNb, s, filteredTweets)
}

func displayStatus(lineNb int, s *traffic.Status, tweets []twitter.Tweet) {
	fmt.Printf("%sCommuting traffic for line %d %s\n\n", emoji.Sprint(":train:"), aurora.BrightBlue(lineNb), emoji.Sprint(":train:"))
	fmt.Printf("\t%s\n\n", s)
	if tweets != nil {
		for _, tweet := range tweets {
			fmt.Println(tweet.Render())
		}
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.commuting-traffic-info.yaml)")
	rootCmd.PersistentFlags().StringVarP(&formatType, "format", "f", "console", "format output \npossible values: \"console\", \"json\", \"rocket.chat\"")
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

func getRocketChatWebhook() string {
	webhook := config.GetRocketChatWebhook()
	if webhook == "" {
		log.Println("No Rocket.Chat webhook configured. Initializing it")
		config.InitRocketChatWebhook(cfgFile)
		webhook = config.GetRocketChatWebhook()
	}
	return webhook
}
