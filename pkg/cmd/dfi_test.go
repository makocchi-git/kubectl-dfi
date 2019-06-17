package cmd

import (
	"bytes"
	"os"
	"reflect"
	"strings"
	"testing"

	color "github.com/gookit/color"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/makocchi-git/kubectl-dfi/pkg/table"
)

// test node object
var testNodes = []v1.Node{
	{
		ObjectMeta: metav1.ObjectMeta{Name: "node1"},
		Status: v1.NodeStatus{
			Images: []v1.ContainerImage{
				{
					Names:     []string{"image1", "image2"},
					SizeBytes: 1000,
				},
			},
			Capacity: v1.ResourceList{
				v1.ResourceEphemeralStorage: *resource.NewQuantity(10*1000*1000*1000, resource.DecimalSI),
			},
			Allocatable: v1.ResourceList{
				v1.ResourceEphemeralStorage: *resource.NewQuantity(5*1000*1000*1000, resource.DecimalSI),
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
			Capacity: v1.ResourceList{
				v1.ResourceEphemeralStorage: *resource.NewQuantity(10*1000*1000*1000, resource.DecimalSI),
			},
			Allocatable: v1.ResourceList{
				v1.ResourceEphemeralStorage: *resource.NewQuantity(5*1000*1000*1000, resource.DecimalSI),
			},
		},
	},
}

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

func TestNewCmdDfi(t *testing.T) {

	rootCmd := NewCmdDfi(
		genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
		"v0.0.1",
		"abcd123",
		"1234567890",
	)

	// Version check
	t.Run("version", func(t *testing.T) {
		expected := "Version: v0.0.1, GitCommit: abcd123, BuildDate: 1234567890\n"
		actual, err := executeCommand(rootCmd, "--version")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// if strings.Contains(output, "kubectl-dfi")
		if actual != expected {
			t.Errorf("expected(%s) differ (got: %s)", expected, actual)
			return
		}
	})

	// Usage
	t.Run("usage", func(t *testing.T) {
		expected := "kubectl dfi [flags]"
		actual, err := executeCommand(rootCmd, "--help")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !strings.Contains(actual, expected) {
			t.Errorf("expected(%s) differ (got: %s)", expected, actual)
			return
		}
	})

	// Unknown option
	t.Run("usage", func(t *testing.T) {
		expected := "unknown flag: --very-very-bad-option"
		_, err := executeCommand(rootCmd, "--very-very-bad-option")
		if err == nil {
			t.Errorf("unexpected error: should return exit")
			return
		}

		if err.Error() != expected {
			t.Errorf("expected(%s) differ (got: %s)", expected, err.Error())
			return
		}
	})

	// RunE Validation error
	t.Run("RunE validation error", func(t *testing.T) {

		c := NewCmdDfi(
			genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
			"v0.0.1",
			"abcd123",
			"1234567890",
		)

		err := c.ParseFlags([]string{"--crit-threshold=5"})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		rerr := c.RunE(c, []string{})
		if rerr == nil {
			t.Errorf("unexpected error: should return error")
			return
		}

		expected := "can not set critical threshold less than warn threshold (warn:25 crit:5)"
		if rerr.Error() != expected {
			t.Errorf("expected(%s) differ (got: %s)", expected, rerr.Error())
			return
		}
	})
}

func TestComplete(t *testing.T) {

	cmd := &cobra.Command{
		Use:   "test",
		Short: "test short",
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

func TestRun(t *testing.T) {

}

func TestDfi(t *testing.T) {

	t.Run("without image count", func(t *testing.T) {

		buffer := &bytes.Buffer{}
		o := &DfiOptions{
			nocolor: true,
			table: table.NewOutputTable(buffer),
		}

		// ---
		// NAME    IMAGE USED   ALLOCATABLE   CAPACITY    %USED
		// node1   1K           5000000K      10000000K   0%
		// node2   2K           5000000K      10000000K   0%
		// ---
		lines := []string{
			"NAME    IMAGE USED   ALLOCATABLE   CAPACITY    %USED",
			"node1   1K           5000000K      10000000K   0%",
			"node2   2K           5000000K      10000000K   0%",
			"",
		}
		expected := strings.Join(lines, "\n")
		if err := o.dfi(testNodes); err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		if buffer.String() != expected {
			t.Errorf("expected(%s) differ (got: %s)", expected, buffer.String())
		}
	})

	t.Run("with image count", func(t *testing.T) {

		buffer := &bytes.Buffer{}
		o := &DfiOptions{
			nocolor: true,
			table: table.NewOutputTable(buffer),
			count: true,
		}

		// ---
		// NAME    IMAGE USED   ALLOCATABLE   CAPACITY    %USED
		// node1   1K(1)        5000000K      10000000K   0%
		// node2   2K(1)        5000000K      10000000K   0%
		// ---
		lines := []string{
			"NAME    IMAGE USED   ALLOCATABLE   CAPACITY    %USED",
			"node1   1K(1)        5000000K      10000000K   0%",
			"node2   2K(1)        5000000K      10000000K   0%",
			"",
		}
		expected := strings.Join(lines, "\n")
		if err := o.dfi(testNodes); err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		if buffer.String() != expected {
			t.Errorf("expected(%s) differ (got: %s)", expected, buffer.String())
		}
	})
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

	if err := o.listImagesOnNode(testNodes); err != nil {
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

// Test Helper
func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOutput(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}
