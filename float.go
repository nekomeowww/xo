package xo

import (
	"regexp"
	"strconv"
)

func IsDecimalsPlacesValid(num float64, decimalPlaces int) bool {
	regex := `^(([1-9]\d*)|(0))(\.\d{0,` + strconv.Itoa(decimalPlaces) + `})?$`
	ok := regexp.MustCompile(regex).MatchString(strconv.FormatFloat(num, 'f', -1, 64))

	return ok
}
