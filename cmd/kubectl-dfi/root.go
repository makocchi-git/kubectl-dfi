package main

import (
	"os"

	"github.com/makocchi-git/kubectl-dfi/pkg/cmd"
	"github.com/spf13/pflag"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	// for goreleaser
	// https://goreleaser.com/customization/#Builds
	version = "master"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	flags := pflag.NewFlagSet("kubectl-dfi", pflag.ExitOnError)
	pflag.CommandLine = flags

	root := cmd.NewCmdDf(
		genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
		version,
		commit,
		date,
	)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
