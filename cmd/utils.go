package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func urlToFileName(url string) string {
	// Replace protocol handler so it just ends up being one underscore
	url = strings.Replace(url, "://", "_", 1)

	// Replace characters not matching the below regex with underscores
	fileName := make([]string, 0)
	re := regexp.MustCompile(`([a-zA-Z0-9])`)
	for _, char := range strings.Split(url, "") {
		if re.MatchString(char) {
			fileName = append(fileName, char)
		} else {
			fileName = append(fileName, "_")
		}
	}

	// Add extension
	return fmt.Sprintf("%s.html", strings.Join(fileName, ""))
}

func writeFile(fileName string, fileContents string) {
	err := ioutil.WriteFile(fileName, []byte(fileContents), 0600)
	if err != nil {
		errMsg(fmt.Sprintf("Error writing to file %s", fileName))
	}
	infoMsg(fmt.Sprintf("Writing to file %s", fileName))
}

func openBrowser(fileName string) {
	infoMsg(fmt.Sprintf("Attempting to use browser: %s", browserPath))

	// Build command string and try to run
	cmd := exec.Command(browserPath, fileName)
	err := cmd.Run()
	if err != nil {
		errMsg(fmt.Sprintf("Error running browser program: %s", browserPath))
	}
}

func errMsg(errorMessage string) {
	fmt.Println(fmt.Sprintf("[!] %s", errorMessage))
	os.Exit(1)
}

func infoMsg(infoMsg string) {
	if verbose {
		fmt.Println(fmt.Sprintf("[*] %s", infoMsg))
	}
}
