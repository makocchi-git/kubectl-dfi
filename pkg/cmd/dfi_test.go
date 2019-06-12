package cmd

import (
	"testing"
)

func TestValidate(t *testing.T) {

	var tests = []struct {
		description string
		warn        int64
		crit        int64
		expected    string
	}{
		{"crit < warn", 25, 10, "can not set critical threshold less than warn threshold (warn:25 crit:10)"},
		{"crit > warn", 25, 30, ""},
		{"warn = crit", 25, 25, ""},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			o := &DfiOptions{
				warnThreshold: test.warn,
				critThreshold: test.crit,
			}
			actual := o.Validate()
			if actual != nil && actual.Error() != test.expected {
				t.Errorf(
					"[%s] expected(%#v) differ (got: %#v)",
					test.description,
					test.expected,
					actual,
				)
				return
			}
		})
	}

}

func TestToUnit(t *testing.T) {

	var tests = []struct {
		description string
		input       int64
		binPrefix   bool
		withoutunit bool
		expected    string
	}{
		{"si prefix without unit", 12345, false, true, "12"},
		{"si prefix with unit", 6000, false, false, "6K"},
		{"binary prefix without unit", 12345, true, true, "12"},
		{"binary prefix with unit", 6000, true, false, "5Ki"},
		{"0 case", 0, true, false, "N/A"},
	}

	o := &DfiOptions{
		bytes: false,
		kByte: true,
		mByte: false,
		gByte: false,
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			o.withoutUnit = test.withoutunit
			o.binPrefix = test.binPrefix
			actual := o.toUnit(test.input)
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

func TestGetImageDiskUsage(t *testing.T) {

	var tests = []struct {
		description string
		used        int64
		capacity    int64
		nocolor     bool
		expected    string
	}{
		{"10%", 10, 100, false, "10%"},
		{"100%", 100, 100, false, "100%"},
		{"over 100%", 123, 100, false, "100%"},
		{"N/A", 0, 0, false, "N/A"},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			o := &DfiOptions{
				nocolor:       test.nocolor,
				warnThreshold: 25,
				critThreshold: 50,
			}
			actual := o.getImageDiskUsage(test.used, test.capacity)
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
