package utils

import (
	"regexp"
)

// (年)-(月)-(天)
var dateRegex string = `(\d{4})-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])`

// (時)-(分)-(秒)
var timeRegex string = `([01][0-9]|2[0123]):([0-5][0-9]):([0-5][0-9])`

func DateTimeRegex(datetime string) bool {
	re := regexp.MustCompile(dateRegex + ` ` + timeRegex)
	return re.MatchString(datetime)
}

func DateRegex(date string) bool {
	re := regexp.MustCompile(dateRegex)
	return re.MatchString(date)
}

func TimeRegex(time string) bool {
	re := regexp.MustCompile(timeRegex)
	return re.MatchString(time)
}
