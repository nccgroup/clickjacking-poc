package cmd

import (
	"fmt"
	"strings"
	"github.com/spf13/cobra"
	"os"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	url string
	outputToFile bool
	browserFString string
	noPrint bool
	verbose bool

	title string
	bodyStyle string
	headerStyle string
	headerMessage string
	iframeStyle string
)

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

		// Get vars. Priority being (highlest to lowest): cli arg, config, defaults
		title = viper.GetString("title")
		bodyStyle = viper.GetString("body-style")
		headerStyle = viper.GetString("header-style")
		headerMessage = viper.GetString("header-message")
		iframeStyle = viper.GetString("iframe-style")

		// Escape " in url
		url = strings.Replace(url, "\"", "\\\"", -1)

		// Build template
		html := fmt.Sprintf(template, title, bodyStyle, headerStyle, headerMessage, url, iframeStyle)

		// Output to file
		outputToFile = viper.GetBool("output-to-file")
		fileName := urlToFileName(url)
		if outputToFile {
			writeFile(fileName, html)
		}

		// Open in browser
		browserFString = viper.GetString("browser-fstring")
		if browserFString != "" && outputToFile == false {
			errMsg("Can't set output to file as false when trying to open file in browser")
		} else {
			openBrowser(fileName, browserFString)
		}

		// Output to stdout if not suppressed
		if ! noPrint {
			fmt.Println(html)
		}
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

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.clickjacking-poc.yaml)")

	rootCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "URL to make a proof of concept for (required)")
	rootCmd.MarkFlagRequired("url")

	// Optional input
	//TODO make this flag editable
	rootCmd.PersistentFlags().BoolVarP(&outputToFile, "output-to-file", "f", true, "Select this if you want to write a file of HTML to the current dir")
	viper.BindPFlag("output-to-file", rootCmd.PersistentFlags().Lookup("output-to-file"))

	rootCmd.PersistentFlags().StringVarP(&browserFString, "browser-fstring", "p", "", "Format string to launch browser from the command line (%s marks file name in current directory")
	viper.BindPFlag("browser-fstring", rootCmd.PersistentFlags().Lookup("browser-fstring"))

	rootCmd.PersistentFlags().BoolVarP(&noPrint, "no-print", "n", false, "Use this flag to suppress html output")
	viper.BindPFlag("no-print", rootCmd.PersistentFlags().Lookup("no-print"))

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable information messages")

	// Optional formatting args
	titleDefault := "Framed Web Application"
	rootCmd.PersistentFlags().StringVarP(&title, "title", "t", titleDefault, "Title of the PoC page")
	viper.BindPFlag("title", rootCmd.PersistentFlags().Lookup("title"))

	bodyStyleDefault := "background-color:black"
	rootCmd.PersistentFlags().StringVarP(&bodyStyle, "body-style", "b", bodyStyleDefault, "CSS style to be applied to the body")
	viper.BindPFlag("body-style", rootCmd.PersistentFlags().Lookup("body-style"))

	headerStyleDefault := "color:white;"
	rootCmd.PersistentFlags().StringVarP(&headerStyle, "header-style", "s", headerStyleDefault, "CSS style to be applied to the header")
	viper.BindPFlag("header-style", rootCmd.PersistentFlags().Lookup("header-style"))

	headerMessageDefault := "The following shows the application embedded in a third party page:"
	rootCmd.PersistentFlags().StringVarP(&headerMessage, "header-message", "m", headerMessageDefault, "Header message above the ifrome")
	viper.BindPFlag("header-message", rootCmd.PersistentFlags().Lookup("header-message"))

	iframeStyleDefault := "width:90%%;height:90%%"
	rootCmd.PersistentFlags().StringVarP(&iframeStyle, "iframe-style", "i", iframeStyleDefault, "CSS style of the iframe")
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
		if verbose {
			fmt.Println("[*] Using config file:", viper.ConfigFileUsed())
		}
	}
}
