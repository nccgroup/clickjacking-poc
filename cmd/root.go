package cmd

import (
	"fmt"
	"strings"
	"github.com/spf13/cobra"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var url string
var title string
var bodyStyle string
var headerStyle string
var headerMessage string
var iframeStyle string

var template string = `<html>
<head><title>%s</title></head>
<body style="%s">
<br />
<h3 style="%s">%s</h3>
<iframe src="%s" style="%s"></iframe>
</body>
</html`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "clickjacking-poc",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Get vars. Priority being (highlest to lowest)
		// CLI
		// Config
		// Defaults
		title = viper.GetString("title")
		bodyStyle = viper.GetString("body-style")
		headerStyle = viper.GetString("header-style")
		headerMessage = viper.GetString("header-message")
		iframeStyle = viper.GetString("iframe-style")

		// Escape " in url
		url = strings.Replace(url, "\"", "\\\"", -1)

		// Build and print template
		html := fmt.Sprintf(template, title, bodyStyle, headerStyle, headerMessage, url, iframeStyle)

		//TODO add flag to writing to file
		fmt.Println(urlToFilename)

		//TODO add file for automatically launching browser
		fmt.Println(html)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.clickjacking-poc.yaml)")

	rootCmd.PersistentFlags().StringVar(&url, "url", "", "URL to make a proof of concept for (required)")
	rootCmd.MarkFlagRequired("url")

	// Optional args
	titleDefault := "Framed Web Application"
	rootCmd.PersistentFlags().StringVarP(&title, "title", "t", titleDefault, "Title of the PoC page")
	viper.BindPFlag("title", rootCmd.PersistentFlags().Lookup("title"))

	bodyStyleDefault := "background-color:black"
	rootCmd.PersistentFlags().StringVar(&bodyStyle, "body-style", bodyStyleDefault, "CSS style to be applied to the body")
	viper.BindPFlag("body-style", rootCmd.PersistentFlags().Lookup("body-style"))

	headerStyleDefault := "color:white;"
	rootCmd.PersistentFlags().StringVar(&headerStyle, "header-style", headerStyleDefault, "CSS style to be applied to the header")
	viper.BindPFlag("header-style", rootCmd.PersistentFlags().Lookup("header-style"))

	headerMessageDefault := "The following shows the application embedded in a third party page:"
	rootCmd.PersistentFlags().StringVar(&headerMessage, "header-message", headerMessageDefault, "Header message above the ifrome")
	viper.BindPFlag("header-message", rootCmd.PersistentFlags().Lookup("header-message"))

	iframeStyleDefault := "width:90%%;height:90%%"
	rootCmd.PersistentFlags().StringVar(&iframeStyle, "iframe-style", iframeStyleDefault, "CSS style of the iframe")
	viper.BindPFlag("iframe-style", rootCmd.PersistentFlags().Lookup("iframe-style"))

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

		// Search config in home directory with name ".clickjacking-poc" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".clickjacking-poc")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//TODO add verbosity flag for this
		fmt.Println("[*] Using config file:", viper.ConfigFileUsed())
	}
}
