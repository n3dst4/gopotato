package cmd

import (
	"fmt"
	"regexp"
)

var longDesc = `
Gopotato is a journal manager for Go. It is a rewrite of Potato.
`

var journalFilenameRegexString = `^(\d{4})[-_](\d{2})[-_](\d{2}).(?:md|markdown)$`

var journalFilenameRegex *regexp.Regexp

func init() {
	var err error
	journalFilenameRegex, err = regexp.Compile(journalFilenameRegexString)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return
	}
}
