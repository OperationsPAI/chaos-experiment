package main

import (
	"os"

	"github.com/OperationsPAI/chaos-experiment/pkg/chaoscli"
)

func main() {
	cmd := chaoscli.NewRootCmd(chaoscli.NewDirectSubmitter(nil, os.Stdout))
	cmd.Use = "chaos-exp"
	cmd.Short = "Standalone chaos CLI that submits directly to Chaos Mesh"
	cmd.PersistentFlags().Bool("dry-run", false, "Print the generated fault spec instead of applying it")
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
