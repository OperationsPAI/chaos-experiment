package chaoscli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

type captureSubmitter struct {
	specs []Spec
}

func (c *captureSubmitter) Submit(_ context.Context, spec Spec) error {
	c.specs = append(c.specs, spec)
	return nil
}

func TestRootContainsMVPSubcommands(t *testing.T) {
	root := NewRootCmd(NewPrintSubmitter(&bytes.Buffer{}))

	paths := walkCommands(root)
	expected := []string{
		"chaos network delay",
		"chaos network loss",
		"chaos http abort",
		"chaos jvm latency",
		"chaos stress cpu",
		"chaos pod failure",
	}

	for _, want := range expected {
		found := false
		for _, path := range paths {
			if path == want {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("missing command path %q in %v", want, paths)
		}
	}
}

func TestNetworkDelayBuildsSpec(t *testing.T) {
	submitter := &captureSubmitter{}
	root := NewRootCmd(submitter)
	root.SetArgs([]string{
		"network", "delay",
		"--namespace", "exp",
		"--app", "frontend",
		"--target-service", "checkout",
		"--duration", "2m",
		"--latency", "120",
		"--correlation", "50",
		"--jitter", "10",
		"--direction", "both",
	})

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if len(submitter.specs) != 1 {
		t.Fatalf("expected one submitted spec, got %d", len(submitter.specs))
	}
	got := submitter.specs[0]
	if got.Type != "NetworkDelay" || got.Namespace != "exp" || got.Target != "frontend" || got.Duration != "2m" {
		t.Fatalf("unexpected spec: %+v", got)
	}
	if got.Params["target_service"] != "checkout" {
		t.Fatalf("expected target_service=checkout, got %#v", got.Params["target_service"])
	}
}

func walkCommands(cmd *cobra.Command) []string {
	var paths []string
	var visit func(*cobra.Command, []string)
	visit = func(current *cobra.Command, prefix []string) {
		parts := append(prefix, current.Name())
		paths = append(paths, strings.Join(parts, " "))
		for _, child := range current.Commands() {
			if !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() {
				continue
			}
			visit(child, parts)
		}
	}
	visit(cmd, nil)
	return paths
}
