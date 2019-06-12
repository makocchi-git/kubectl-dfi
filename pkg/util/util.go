package util

import (
	"strings"

	color "github.com/gookit/color"
	"github.com/makocchi-git/kubectl-dfi/pkg/constants"
)

// GetSiUnit defines unit for usage (SI prefix)
// If multiple options are selected, returns a biggest unit
func GetSiUnit(b, k, m, g bool) (int64, string) {

	if g {
		return constants.UnitGigaBytes, constants.UnitGigaBytesStr
	}

	if m {
		return constants.UnitMegaBytes, constants.UnitMegaBytesStr
	}

	if k {
		return constants.UnitKiloBytes, constants.UnitKiloBytesStr
	}

	if b {
		return constants.UnitBytes, constants.UnitBytesStr
	}

	// default output is "kilobytes"
	return constants.UnitKiloBytes, constants.UnitKiloBytesStr
}

// GetBinUnit defines unit for usage (Binary prefix)
// If multiple options are selected, returns a biggest unit
func GetBinUnit(b, k, m, g bool) (int64, string) {

	if g {
		return constants.UnitGibiBytes, constants.UnitGibiBytesStr
	}

	if m {
		return constants.UnitMibiBytes, constants.UnitMibiBytesStr
	}

	if k {
		return constants.UnitKibiBytes, constants.UnitKibiBytesStr
	}

	if b {
		return constants.UnitBytes, constants.UnitBytesStr
	}

	// default output is "kibibytes"
	return constants.UnitKibiBytes, constants.UnitKibiBytesStr
}

// JoinTab joins string slice with tab
func JoinTab(s []string) string {
	return strings.Join(s, "\t")
}

// SetPercentageColor returns colored string
//        percentage < warn : Green
// warn < percentage < crit : Yellow
// crit < percentage        : Red
func SetPercentageColor(s *string, p, warn, crit int64) {
	green := color.FgGreen.Render
	yellow := color.FgYellow.Render
	red := color.FgRed.Render

	if p < warn {
		*s = green(*s)
		return
	}

	if p < crit {
		*s = yellow(*s)
		return
	}

	*s = red(*s)
}

// ColorImageTag is colorize image tag
func ColorImageTag(image *string) {

	if strings.Contains(*image, ":") {

		s := strings.Split(*image, ":")
		sl := len(s)

		// color last one
		yellow := color.FgYellow.Render
		tag := yellow(s[sl-1])

		// rejoin strings
		*image = strings.Join(s[0:sl-1], ":") + ":" + tag
	}
}