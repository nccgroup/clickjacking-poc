package cmd

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"os"
	"os/exec"
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

	// Add extention
	return fmt.Sprintf("%s.html", strings.Join(fileName, ""))
}

func writeFile(fileName string, fileContents string) {
	err := ioutil.WriteFile(fileName, []byte(fileContents), 0600)
	if err != nil {
		errMsg(fmt.Sprintf("Error writing to file %s", fileName))
	}
	if verbose {
		infoMsg(fmt.Sprintf("Writing to file %s", fileName))
	}
}

func openBrowser(fileName string, browserFString string) {
	// Ensure that the correct number of %s exist in browserFString
	if strings.Count(browserFString, "%s") == 0 {
		errMsg("%s not supplied in browser format string!")
	}

	// Build command string and try to run
	command := fmt.Sprintf(browserFString, fileName)
	if verbose {
		infoMsg(fmt.Sprintf("Attempting to run command: %s", command))
	}
	cmd := exec.Command(command)
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	err := cmd.Run()
	if err != nil {
		panic(err)
		errMsg(fmt.Sprintf("Error running browser program: %s", command))
	}
}

func errMsg(errorMessage string) {
	fmt.Println(fmt.Sprintf("[!] %s", errorMessage))
	os.Exit(1)
}

func infoMsg(infoMsg string) {
	fmt.Println(fmt.Sprintf("[*] %s", infoMsg))
}
