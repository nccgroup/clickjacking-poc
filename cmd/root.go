package cmd

import (
	"fmt"
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
		html := fmt.Sprintf(template, title, bodyStyle, headerStyle, headerMessage, url, iframeStyle)
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
	rootCmd.PersistentFlags().StringVar(&title, "title", titleDefault, "Title of the PoC page")
	viper.BindPFlag("title", rootCmd.PersistentFlags().Lookup("title"))
	viper.SetDefault("title", titleDefault)

	bodyStyleDefault := "background-color:black"
	rootCmd.PersistentFlags().StringVar(&bodyStyle, "body-style", bodyStyleDefault, "CSS style to be applied to the body")
	viper.BindPFlag("body-style", rootCmd.PersistentFlags().Lookup("body-style"))
	viper.SetDefault("body-style", bodyStyleDefault)

	headerStyleDefault := "color:white;"
	rootCmd.PersistentFlags().StringVar(&headerStyle, "header-style", headerStyleDefault, "CSS style to be applied to the header")
	viper.BindPFlag("header-style", rootCmd.PersistentFlags().Lookup("header-style"))
	viper.SetDefault("header-style", headerStyleDefault)

	headerMessageDefault := "The following shows the application embedded in a third party page:"
	rootCmd.PersistentFlags().StringVar(&headerMessage, "header-message", headerMessageDefault, "Header message above the ifrome")
	viper.BindPFlag("header-message", rootCmd.PersistentFlags().Lookup("header-message"))
	viper.SetDefault("header-message", headerMessageDefault)

	iframeStyleDefault := "width:90%%;height:90%%"
	rootCmd.PersistentFlags().StringVar(&iframeStyle, "iframe-style", iframeStyleDefault, "CSS style of the iframe")
	viper.BindPFlag("iframe-style", rootCmd.PersistentFlags().Lookup("iframe-style"))
	viper.SetDefault("iframe-style", iframeStyleDefault)

	// Bind with viper
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))

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
		fmt.Println("[*] Using config file:", viper.ConfigFileUsed())
	}
}
