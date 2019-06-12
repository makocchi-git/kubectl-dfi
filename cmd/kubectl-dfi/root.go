package main

import (
	"flag"
	"os"

	"github.com/makocchi-git/kubectl-dfi/pkg/cmd"
	"github.com/spf13/pflag"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func init() {
	// Initialize glog flags
	if err := flag.CommandLine.Set("all", "false"); err != nil {
		os.Exit(1)
	}
}

func main() {
	flags := pflag.NewFlagSet("kubectl-dfi", pflag.ExitOnError)
	pflag.CommandLine = flags

	root := cmd.NewCmdDf(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
