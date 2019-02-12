package helpanto

import (
	"strings"
	"unicode"
)

func GetRegionFromDC(dc string) string {

	region := strings.TrimFunc(dc, func(r rune) bool {
		return !unicode.IsLetter(r)
	})
	return region

}
