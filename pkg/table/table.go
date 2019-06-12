package table

import (
	"fmt"
	"io"
	"os"

	"github.com/makocchi-git/kubectl-dfi/pkg/util"

	"k8s.io/kubernetes/pkg/printers"
)

// OutputTable is struct of tables for outputs
type OutputTable struct {
	Header string
	Rows   []string
	Output io.Writer
}

// NewOutputTable is an instance of OutputTable
func NewOutputTable() *OutputTable {
	return &OutputTable{
		Output: os.Stdout,
	}
}

// Print shows table output
func (t *OutputTable) Print() {

	// get printer
	printer := printers.GetNewTabWriter(t.Output)

	// write header
	fmt.Fprintln(printer, t.Header)

	// write rows
	for _, row := range t.Rows {
		fmt.Fprintln(printer, row)
	}

	// finish
	printer.Flush()
}

// AddHeader adds row to table
func (t *OutputTable) AddHeader(s []string) {
	t.Header = util.JoinTab(s)
}

// AddRow adds row to table
func (t *OutputTable) AddRow(s []string) {
	t.Rows = append(t.Rows, util.JoinTab(s))
}
