package utils

import "regexp"

//ParseInstance - Helper function to extract pandora instance from filename
func ParseInstance(toParse string) (rs string) {

	var rgx = regexp.MustCompile(`\((.*?)\)`)

	rss := rgx.FindStringSubmatch(toParse)

	rs = rss[1]

	return

}
