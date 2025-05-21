package utils

import (
	"regexp"
	"strings"
)

// TrimMessage trims an "s" string and removes bad chars, links and pings
func TrimMessage(s string) string {
	var not, ping, link *regexp.Regexp
	var err error

	not, err = regexp.Compile("[^a-zA-Z0-9éèêàùûç,;:!.?]")
	if err != nil {
		SendAlert("strings.go - Impossible to compile regex 'not'", err.Error())
		return ""
	}
	ping, err = regexp.Compile("<(@&?|#)[0-9]{18}>")
	if err != nil {
		SendAlert("strings.go - Impossible to compile regex 'ping'", err.Error())
		return ""
	}
	link, err = regexp.Compile("https?://[a-zA-Z0-9.]+[.][a-z]+.*")
	if err != nil {
		SendAlert("strings.go - Impossible to compile regex 'link'", err.Error())
		return ""
	}

	s = ping.ReplaceAllLiteralString(s, "")
	s = link.ReplaceAllLiteralString(s, "")
	s = not.ReplaceAllLiteralString(s, "")

	return strings.Trim(s, " ")
}
