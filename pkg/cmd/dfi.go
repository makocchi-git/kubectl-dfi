package cmd

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/makocchi-git/kubectl-dfi/pkg/table"
	"github.com/makocchi-git/kubectl-dfi/pkg/util"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/kubectl/util/templates"

	// Initialize all known client auth plugins.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var (
	// DfLong defines long description
	DfLong = templates.LongDesc(`
		Show disk resources of images on Kubernetes nodes.
	`)

	// DfExample defines command examples
	DfExample = templates.Examples(`
		# Show image usage of Kubernetes nodes.
		kubectl dfi

		# Using label selector.
		kubectl dfi -l key=value

		# Use image count with image disk usage.
		kubectl dfi --count

		# Print raw(bytes) usage.
		kubectl dfi --bytes --without-unit

		# Using binary prefix unit (GiB, MiB, etc)
		kubectl dfi -g -B

		# List images on nodes.
		kubectl dfi --list
	`)
)

// DfOptions is struct of df options
type DfOptions struct {
	configFlags *genericclioptions.ConfigFlags
	genericclioptions.IOStreams

	// general options
	labelSelector string
	count         bool
	table         *table.OutputTable

	// unit options
	bytes       bool
	kByte       bool
	mByte       bool
	gByte       bool
	withoutUnit bool
	binPrefix   bool

	// color output options
	nocolor       bool
	warnThreshold int64
	critThreshold int64

	// list options
	list bool
}

// NewDfOptions is an instance of DfOptions
func NewDfOptions(streams genericclioptions.IOStreams) *DfOptions {
	return &DfOptions{
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
		table:         table.NewOutputTable(),
	}
}

// NewCmdDf is a cobra command wrapping
func NewCmdDf(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewDfOptions(streams)

	cmd := &cobra.Command{
		Use:     fmt.Sprintf("kubectl dfi"),
		Short:   "Show disk resources of images on Kubernetes nodes.",
		Long:    DfLong,
		Example: DfExample,
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(c, args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			c.SilenceUsage = true
			if err := o.Run(args); err != nil {
				return err
			}

			return nil
		},
	}

	// bool options
	cmd.Flags().BoolVarP(&o.bytes, "bytes", "b", o.bytes, `Use 1-byte (1-Byte) blocks rather than the default.`)
	cmd.Flags().BoolVarP(&o.kByte, "kilobytes", "k", o.kByte, `Use 1024-byte (1-Kbyte) blocks rather than the default.`)
	cmd.Flags().BoolVarP(&o.mByte, "megabytes", "m", o.mByte, `Use 1048576-byte (1-Mbyte) blocks rather than the default.`)
	cmd.Flags().BoolVarP(&o.gByte, "gigabytes", "g", o.gByte, `Use 1073741824-byte (1-Gbyte) blocks rather than the default.`)
	cmd.Flags().BoolVarP(&o.binPrefix, "binary-prefix", "B", o.binPrefix, `Use 1024 for basic unit calculation instead of 1000. (print like "KiB")`)
	cmd.Flags().BoolVarP(&o.withoutUnit, "without-unit", "", o.withoutUnit, `Do not print size with unit string.`)
	cmd.Flags().BoolVarP(&o.count, "count", "c", o.count, `Print number of images.`)
	cmd.Flags().BoolVarP(&o.nocolor, "no-color", "", o.nocolor, `Print without ansi color.`)
	cmd.Flags().BoolVarP(&o.list, "list", "", o.list, `Show image list on node.`)

	// int64 options
	cmd.Flags().Int64VarP(&o.warnThreshold, "warn-threshold", "", o.warnThreshold, `Threshold of warn(yellow) color for USED column.`)
	cmd.Flags().Int64VarP(&o.critThreshold, "crit-threshold", "", o.critThreshold, `Threshold of critical(red) color for USED column.`)

	// string option
	cmd.Flags().StringVarP(&o.labelSelector, "selector", "l", o.labelSelector, `Selector (label query) to filter on.`)

	o.configFlags.AddFlags(cmd.Flags())

	// add the klog flags
	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	return cmd
}

// Complete sets all information required for opening the service
func (o *DfOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

// Validate ensures that all required arguments and flag values are provided
func (o *DfOptions) Validate() error {

	if o.warnThreshold > o.critThreshold {
		return fmt.Errorf(
			"Can not set critical threshold less than warn threshold (warn:%d crit:%d)", o.warnThreshold, o.critThreshold,
		)
	}

	return nil
}

// Run opens the service in the browser
func (o *DfOptions) Run(args []string) error {

	// get k8s client
	restConfig, err := o.configFlags.ToRESTConfig()
	if err != nil {
		return err
	}
	client := kubernetes.NewForConfigOrDie(restConfig)

	// get nodes
	nodes := []v1.Node{}
	if len(args) > 0 {
		for _, a := range args {
			n, nerr := client.CoreV1().Nodes().Get(a, metav1.GetOptions{})
			if nerr != nil {
				return fmt.Errorf("Failed to get node: %v", nerr)
			}
			nodes = append(nodes, *n)
		}
	} else {
		na, naerr := client.CoreV1().Nodes().List(metav1.ListOptions{LabelSelector: o.labelSelector})
		if naerr != nil {
			return fmt.Errorf("Failed to get nodes: %v", naerr)
		}
		nodes = append(nodes, na.Items...)
	}

	// list images and return
	if o.list {
		if err := o.listImagesOnNode(nodes); err != nil {
			return err
		}
		return nil
	}

	// print df image usage
	if err := o.dfi(nodes); err != nil {
		return err
	}

	return nil
}

// toUnit caluculate and add unit for int64
func (o *DfOptions) toUnit(i int64) string {

	var unitbytes int64
	var unitstr string

	if o.binPrefix {
		unitbytes, unitstr = util.GetBinUnit(o.bytes, o.kByte, o.mByte, o.gByte)
	} else {
		unitbytes, unitstr = util.GetSiUnit(o.bytes, o.kByte, o.mByte, o.gByte)
	}
	
	// -H adds human readable unit
	unit := ""
	if !o.withoutUnit {
		unit = unitstr
	}

	return strconv.FormatInt(i/unitbytes, 10) + unit
}

// dfi prints image disk usage
func (o *DfOptions) dfi(nodes []v1.Node) error {

	// set printer header
	headers := []string{"NAME", "IMAGE USED", "ALLOCATABLE", "CAPACITY", "%USED"}
	o.table.AddHeader(headers)

	// node loop
	for _, node := range nodes {

		// node name
		name := node.ObjectMeta.Name

		// get status.capacity
		capacity, cerr := node.Status.Capacity.StorageEphemeral().AsInt64()
		if !cerr {
			return fmt.Errorf("Can not get ephemeral storage capacity")
		}

		// get status.allocatable
		allocatable, aerr := node.Status.Allocatable.StorageEphemeral().AsInt64()
		if !aerr {
			return fmt.Errorf("Can not get ephemeral storage capacity")
		}

		// get used storage by images and count images
		used, count := getImageUsage(node.Status.Images)

		// with image count
		icount := ""
		if o.count {
			icount = fmt.Sprintf("(%d)", count)
		}

		// image disk usage
		percentage := (used * 100) / capacity

		// set color
		columnsPercent := strconv.FormatInt(percentage, 10) + "%"
		if !o.nocolor {
			util.SetPercentageColor(&columnsPercent, percentage, o.warnThreshold, o.critThreshold)
		}

		// columns
		columnUsed := o.toUnit(used) + icount
		columnAllocatable := o.toUnit(allocatable)
		columnCapacity := o.toUnit(capacity)

		row := []string{
			name,
			columnUsed,
			columnAllocatable,
			columnCapacity,
			columnsPercent,
		}
		o.table.AddRow(row)
	}

	o.table.Print()

	return nil
}

func (o *DfOptions) listImagesOnNode(nodes []v1.Node) error {

	// set printer header
	headers := []string{"NAME", "IMAGE SIZE", "IMAGE NAME"}
	o.table.AddHeader(headers)

	// node loop
	for _, node := range nodes {

		// node name
		name := node.ObjectMeta.Name

		for _, i := range node.Status.Images {
			imageName := i.Names[0]
			if len(i.Names) > 1 {
				imageName = i.Names[1]
			}

			// color tag
			if !o.nocolor {
				util.ColorImageTag(&imageName)
			}

			row := []string{name, o.toUnit(i.SizeBytes), imageName}
			o.table.AddRow(row)
		}
	}

	o.table.Print()

	return nil
}

// getImageUsage returns total image size and count
func getImageUsage(images []v1.ContainerImage) (int64, int) {

	var s int64
	for _, image := range images {
		s = s + image.SizeBytes
	}
	return s, len(images)
}
