package cmd

import (
	"regexp"
	"fmt"
)

func urlToFilename(url string) string {

	re := regexp.MustCompile(`(?![a-zA-Z-_0-9]+)`)
	fmt.Println(re.ReplaceAllLiteralString(url, "_"))
	return ":p;"
}
