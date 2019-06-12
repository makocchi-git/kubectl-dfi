package cmd

import (
	"testing"

	v1 "k8s.io/api/core/v1"
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

func TestGetImageUsage(t *testing.T) {

	var expectedSize int64
	var expectedCount int

	expectedSize = 11111
	expectedCount = 5
	images := []v1.ContainerImage{
		{Names: []string{"image1"}, SizeBytes: 1},
		{Names: []string{"image2"}, SizeBytes: 10},
		{Names: []string{"image3"}, SizeBytes: 100},
		{Names: []string{"image4"}, SizeBytes: 1000},
		{Names: []string{"image5"}, SizeBytes: 10000},
	}

	actualSize, actualCount := getImageUsage(images)

	if actualSize != expectedSize {
		t.Errorf("Size check: expected(%d) differ (got: %d)", expectedSize, actualSize)
	}

	if actualCount != expectedCount {
		t.Errorf("Count check: expected(%d) differ (got: %d)", expectedCount, actualCount)
	}

}
