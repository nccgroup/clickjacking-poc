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
	browserPath string
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
</html>`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "clickjacking-poc",
	Short: "A tool to generate clickjacking proof of concepts for a given URL",
	Long: fmt.Sprintf(`This tool generates a clickjacking proof of concept for a provided URL.
By default the generated HTML is printed to stdout and also saved to a file based on the name of the site.

Default usage:
clickjacking-poc -u http://wikipedia.com

The tool can also be configured to open the generated file in a browser. This has been tested with firefox and chromium on Linux:
clickjacking-poc -u http://wikipedia.com -b chromium-browser

Flags also exist for suppressing the stdout and file output.

A number of optional flags exist that can be used to inject customized styling into the template.
To see these, check the output of clickjacking-poc -h

The current template for this tool is:
%s

The tool is also built with with viper in mind for using a configuration file to override the tool defaults and lose the requirement of passing the same flags again and again at the command line.

In theory all viper supported formats should be supported. The configus should start with .clickjacking-poc and be placed in your home directory (use --config to change location). Examples can be found on GitHub. For Viper configuration see https://github.com/spf13/viper`, template),
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

		// Open generated file in browser
		browserPath = viper.GetString("browser-path")

		// If the output of the html to a file has been suppresed we can't open
		// it in a browser
		if browserPath != "" && outputToFile == false {
			errMsg("Can't set output to file as false when trying to open file in browser")

		// If browser path not set
		} else if browserPath != ""{
			openBrowser(fileName)
		}

		// Output to stdout if not suppressed
		if noPrint == true {
			infoMsg("Generated HTML:")
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
	rootCmd.PersistentFlags().BoolVarP(&outputToFile, "output-to-file", "f", true, "Toggle whether or not you want to write the HTML to a file (default true) (required for opening in browser)")
	viper.BindPFlag("output-to-file", rootCmd.PersistentFlags().Lookup("output-to-file"))

	rootCmd.PersistentFlags().StringVarP(&browserPath, "browser-path", "b", "", "Path to browser binary (if in PATH only the binary name should be required)")
	viper.BindPFlag("browser-path", rootCmd.PersistentFlags().Lookup("browser-path"))

	rootCmd.PersistentFlags().BoolVarP(&noPrint, "no-print", "n", false, "Use this flag to suppress HTML output")
	viper.BindPFlag("no-print", rootCmd.PersistentFlags().Lookup("no-print"))

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable informational messages")

	// Optional formatting args
	titleDefault := "Framed Web Application"
	rootCmd.PersistentFlags().StringVarP(&title, "title", "t", titleDefault, "Title of the PoC page")
	viper.BindPFlag("title", rootCmd.PersistentFlags().Lookup("title"))

	bodyStyleDefault := "background-color:black"
	rootCmd.PersistentFlags().StringVarP(&bodyStyle, "body-style", "s", bodyStyleDefault, "CSS style to be applied to the body")
	viper.BindPFlag("body-style", rootCmd.PersistentFlags().Lookup("body-style"))

	headerStyleDefault := "color:white;"
	rootCmd.PersistentFlags().StringVarP(&headerStyle, "header-style", "y", headerStyleDefault, "CSS style to be applied to the header")
	viper.BindPFlag("header-style", rootCmd.PersistentFlags().Lookup("header-style"))

	headerMessageDefault := "The following shows the application embedded in a third party page:"
	rootCmd.PersistentFlags().StringVarP(&headerMessage, "header-message", "m", headerMessageDefault, "Message above the iframe")
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
		infoMsg(fmt.Sprintf("Using config file: %s", viper.ConfigFileUsed()))
	} else {
		errMsg(fmt.Sprintf("Error using config file %s ! Is it formatted correctly?", viper.ConfigFileUsed()))
	}
}
