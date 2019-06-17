package cmd

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/makocchi-git/kubectl-dfi/pkg/table"
	"github.com/makocchi-git/kubectl-dfi/pkg/util"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/kubectl/util/templates"
	clientv1 "k8s.io/client-go/kubernetes/typed/core/v1"

	// Initialize all known client auth plugins.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var (
	// DfLong defines long description
	dfiLong = templates.LongDesc(`
		Show disk resources of images on Kubernetes nodes.
	`)

	// DfExample defines command examples
	dfiExample = templates.Examples(`
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

// DfiOptions is struct of df options
type DfiOptions struct {
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

	// k8s node client
	nodeClient   clientv1.NodeInterface
}

// NewDfOptions is an instance of DfOptions
func NewDfiOptions(streams genericclioptions.IOStreams) *DfiOptions {
	return &DfiOptions{
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
}

// NewCmdDf is a cobra command wrapping
func NewCmdDfi(streams genericclioptions.IOStreams, version, commit, date string) *cobra.Command {
	o := NewDfiOptions(streams)

	cmd := &cobra.Command{
		Use:     fmt.Sprintf("kubectl dfi"),
		Short:   "Show disk resources of images on Kubernetes nodes.",
		Long:    dfiLong,
		Example: dfiExample,
		Version: version,
		RunE: func(c *cobra.Command, args []string) error {
			c.SilenceUsage = true

			// TODO: Currently not implemented
			_ = o.Complete(c, args)
			// if err := o.Complete(c, args); err != nil {
			// 	return err
			// }

			if err := o.Validate(); err != nil {
				return err
			}

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

	// version command template
	cmd.SetVersionTemplate("Version: " + version + ", GitCommit: " + commit + ", BuildDate: " + date + "\n")

	return cmd
}

// Complete sets all information required for opening the service
func (o *DfiOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

// Validate ensures that all required arguments and flag values are provided
func (o *DfiOptions) Validate() error {

	if o.warnThreshold > o.critThreshold {
		return fmt.Errorf(
			"can not set critical threshold less than warn threshold (warn:%d crit:%d)", o.warnThreshold, o.critThreshold,
		)
	}

	return nil
}

// Run opens the service in the browser
func (o *DfiOptions) Run(args []string) error {

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
				return fmt.Errorf("failed to get node: %v", nerr)
			}
			nodes = append(nodes, *n)
		}
	} else {
		na, naerr := client.CoreV1().Nodes().List(metav1.ListOptions{LabelSelector: o.labelSelector})
		if naerr != nil {
			return fmt.Errorf("failed to get nodes: %v", naerr)
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

// dfi prints image disk usage
func (o *DfiOptions) dfi(nodes []v1.Node) error {

	// set printer header
	headers := []string{"NAME", "IMAGE USED", "ALLOCATABLE", "CAPACITY", "%USED"}
	o.table.AddHeader(headers)

	// node loop
	for _, node := range nodes {

		// node name
		name := node.ObjectMeta.Name

		// get status.capacity
		capacity, _ := node.Status.Capacity.StorageEphemeral().AsInt64()

		// get status.allocatable
		allocatable, _ := node.Status.Allocatable.StorageEphemeral().AsInt64()

		// get used storage by images and count images
		used, count := util.GetImageUsage(node.Status.Images)

		// with image count
		icount := ""
		if o.count {
			icount = fmt.Sprintf("(%d)", count)
		}

		// columns
		columnsPercent := o.getImageDiskUsage(used, capacity)
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

func (o *DfiOptions) listImagesOnNode(nodes []v1.Node) error {

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

// toUnit calculate and add unit for int64
func (o *DfiOptions) toUnit(i int64) string {

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

	// Old kubernetes do not support capacity attribute.
	if i == 0 {
		return "N/A"
	}

	return strconv.FormatInt(i/unitbytes, 10) + unit
}

func (o *DfiOptions) getImageDiskUsage(used, capacity int64) string {

	var ret string

	// Old kubernetes do not support capacity attribute.
	// So I should avoid panic with "integer divide by zero"
	if capacity == 0 {
		ret = "N/A"
	} else {
		p := (used * 100) / capacity

		// maybe something wrong
		if p > 100 {
			p = 100
		}

		ret = strconv.FormatInt(p, 10) + "%"

		// set color
		if !o.nocolor {
			util.SetPercentageColor(&ret, p, o.warnThreshold, o.critThreshold)
		}
	}

	return ret
}
