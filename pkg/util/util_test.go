package util

import (
	"strconv"
	"testing"

	"github.com/makocchi-git/kubectl-dfi/pkg/constants"

	color "github.com/gookit/color"
)

func TestGetSiUnit(t *testing.T) {
	var tests = []struct {
		description string
		b           bool
		k           bool
		m           bool
		g           bool
		expectedInt int64
		expectedStr string
	}{
		{"all false", false, false, false, false, constants.UnitKiloBytes, constants.UnitKiloBytesStr},
		{"all true", true, true, true, true, constants.UnitGigaBytes, constants.UnitGigaBytesStr},
		{"b only", true, false, false, false, constants.UnitBytes, constants.UnitBytesStr},
		{"g only", false, false, false, true, constants.UnitGigaBytes, constants.UnitGigaBytesStr},
		{"b and k", true, true, false, false, constants.UnitKiloBytes, constants.UnitKiloBytesStr},
		{"k and m", false, true, true, false, constants.UnitMegaBytes, constants.UnitMegaBytesStr},
		{"k and m and g", false, true, true, true, constants.UnitGigaBytes, constants.UnitGigaBytesStr},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actualInt, actualStr := GetSiUnit(test.b, test.k, test.m, test.g)
			if actualInt != test.expectedInt || actualStr != test.expectedStr {
				t.Errorf(
					"[%s] expected(%d, %s) differ (got: %d, %s)",
					test.description,
					test.expectedInt,
					test.expectedStr,
					actualInt,
					actualStr,
				)
				return
			}
		})
	}
}

func TestGetBinUnit(t *testing.T) {
	var tests = []struct {
		description string
		b           bool
		k           bool
		m           bool
		g           bool
		expectedInt int64
		expectedStr string
	}{
		{"all false", false, false, false, false, constants.UnitKibiBytes, constants.UnitKibiBytesStr},
		{"all true", true, true, true, true, constants.UnitGibiBytes, constants.UnitGibiBytesStr},
		{"b only", true, false, false, false, constants.UnitBytes, constants.UnitBytesStr},
		{"g only", false, false, false, true, constants.UnitGibiBytes, constants.UnitGibiBytesStr},
		{"b and k", true, true, false, false, constants.UnitKibiBytes, constants.UnitKibiBytesStr},
		{"k and m", false, true, true, false, constants.UnitMibiBytes, constants.UnitMibiBytesStr},
		{"k and m and g", false, true, true, true, constants.UnitGibiBytes, constants.UnitGibiBytesStr},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actualInt, actualStr := GetBinUnit(test.b, test.k, test.m, test.g)
			if actualInt != test.expectedInt || actualStr != test.expectedStr {
				t.Errorf(
					"[%s] expected(%d, %s) differ (got: %d, %s)",
					test.description,
					test.expectedInt,
					test.expectedStr,
					actualInt,
					actualStr,
				)
				return
			}
		})
	}
}

func TestJoinTab(t *testing.T) {
	var tests = []struct {
		description string
		words       []string
		expected    string
	}{
		{"1 string", []string{"foo"}, "foo"},
		{"2 strings", []string{"foo", "bar"}, "foo\tbar"},
		{"3 strings", []string{"foo", "bar", "baz"}, "foo\tbar\tbaz"},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual := JoinTab(test.words)
			if actual != test.expected {
				t.Errorf(
					"[%s] expected(%s) differ (got: %s)",
					test.description,
					test.expected,
					actual,
				)
				return
			}
		})
	}
}

func TestSetPercentageColor(t *testing.T) {

	var tests = []struct {
		description string
		percentage  int64
		warn        int64
		crit        int64
		expected    string
	}{
		{"5 with warn 10 crit 30", 5, 10, 30, color.Green.Sprint("5%")},
		{"15 with warn 10 crit 30", 15, 10, 30, color.Yellow.Sprint("15%")},
		{"50 with warn 10 crit 30", 50, 10, 30, color.Red.Sprint("50%")},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual := strconv.FormatInt(test.percentage, 10) + "%"
			SetPercentageColor(&actual, test.percentage, test.warn, test.crit)

			if actual != test.expected {
				t.Errorf(
					"[%s] expected(%v) differ (got: %v)",
					test.description,
					actual,
					test.expected,
				)
				return
			}
		})
	}
}

func TestColorImageTag(t *testing.T) {

	yellow := color.FgYellow.Render

	var tests = []struct {
		description string
		image       string
		expected    string
	}{
		{"one delimiter", "abc/def:latest", "abc/def:" + yellow("latest")},
		{"two delimiters", "abc:5000/def:v1.2.3", "abc:5000/def:" + yellow("v1.2.3")},
		{"no delimiters", "abc/def", "abc/def"},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual := test.image
			ColorImageTag(&actual)

			if actual != test.expected {
				t.Errorf(
					"[%s] expected(%v) differ (got: %v)",
					test.description,
					actual,
					test.expected,
				)
				return
			}
		})
	}
}
