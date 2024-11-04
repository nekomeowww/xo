package xo

import (
	"regexp"
	"strings"
)

var rfc1123ValidName = regexp.MustCompile("[^A-Za-z0-9-]")
var rfc1123LeadingTrailingValidName = regexp.MustCompile("^[^A-Za-z0-9]+|[^A-Za-z0-9]+$")

func NormalizeAsRFC1123Name(value string) string {
	// Replace "/" with "-"
	value = strings.ReplaceAll(value, "/", "-")

	// Replace invalid characters with hyphens
	value = rfc1123ValidName.ReplaceAllString(value, "-")

	// Trim hyphens, underscores, and dots from the beginning and end
	value = strings.Trim(value, "-")

	// Ensure it starts and ends with an alphanumeric character
	value = rfc1123LeadingTrailingValidName.ReplaceAllString(value, "")

	// Truncate to maximum length if necessary
	if len(value) > 256 {
		value = value[:256]
	}

	return value
}
