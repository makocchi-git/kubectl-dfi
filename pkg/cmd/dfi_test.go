package cmd

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	color "github.com/gookit/color"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/makocchi-git/kubectl-dfi/pkg/table"
)

func TestNewDfiOptions(t *testing.T) {

	streams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}

	expected := &DfiOptions{
		configFlags:   genericclioptions.NewConfigFlags(true),
		bytes:         false,
		kByte:         false,
		mByte:         false,
		gByte:         false,
		withoutUnit:   false,
		binPrefix:     false,
		count:         false,
		nocolor:       false,
		warnThreshold: 25,
		critThreshold: 50,
		IOStreams:     streams,
		labelSelector: "",
		list:          false,
		table:         table.NewOutputTable(os.Stdout),
	}

	actual := NewDfiOptions(streams)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected(%#v) differ (got: %#v)", expected, actual)
	}
}

func TestComplete(t *testing.T) {

	cmd := &cobra.Command{
		Use:     "test",
		Short:   "test short",
	}

	o := &DfiOptions{}
	if err := o.Complete(cmd, []string{}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

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

func TestListImagesOnNode(t *testing.T) {

	nodes := []v1.Node{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "node1"},
			Status: v1.NodeStatus{
				Images: []v1.ContainerImage{
					{
						Names:     []string{"image1", "image2"},
						SizeBytes: 1000,
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "node2"},
			Status: v1.NodeStatus{
				Images: []v1.ContainerImage{
					{
						Names:     []string{"image1"},
						SizeBytes: 2000,
					},
				},
			},
		},
	}

	buffer := &bytes.Buffer{}
	o := &DfiOptions{
		table: table.NewOutputTable(buffer),
	}

	// ---
	// NAME    IMAGE SIZE   IMAGE NAME
	// node1   1K           image2
	// node2   2K           image1
	// ---
	expected := "NAME    IMAGE SIZE   IMAGE NAME\nnode1   1K           image2\nnode2   2K           image1\n"

	if err := o.listImagesOnNode(nodes); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if buffer.String() != expected {
		t.Errorf("expected(%#v) differ (got: %#v)", expected, buffer.String())
	}

}

func TestGetImageDiskUsage(t *testing.T) {

	red := color.FgRed.Render

	var tests = []struct {
		description string
		used        int64
		capacity    int64
		nocolor     bool
		expected    string
	}{
		{"10%", 10, 100, true, "10%"},
		{"100%", 100, 100, true, "100%"},
		{"over 100%", 123, 100, true, "100%"},
		{"N/A", 0, 0, true, "N/A"},
		{"100% with color", 100, 100, false, red("100%")},
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
