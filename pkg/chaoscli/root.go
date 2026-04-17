package chaoscli

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

type commandOptions struct {
	namespace string
}

type dryRunSubmitter interface {
	DryRun(ctx context.Context, spec Spec, w io.Writer) error
}

func NewRootCmd(s Submitter) *cobra.Command {
	opts := &commandOptions{}

	cmd := &cobra.Command{
		Use:           "chaos",
		Short:         "Submit chaos faults from an imperative CLI",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.PersistentFlags().StringVar(&opts.namespace, "namespace", "ts", "Target namespace")
	cmd.PersistentFlags().Bool("dry-run", false, "Print the generated spec instead of submitting it")

	cmd.AddCommand(
		NewNetworkCmd(s, opts),
		NewHTTPCmd(s, opts),
		NewJVMCmd(s, opts),
		NewStressCmd(s, opts),
		NewPodCmd(s, opts),
	)

	return cmd
}

func submitSpec(cmd *cobra.Command, s Submitter, spec Spec) error {
	if spec.Type == "" {
		return fmt.Errorf("fault type is required")
	}
	if spec.Target == "" {
		return fmt.Errorf("--app is required")
	}
	if spec.Namespace == "" {
		return fmt.Errorf("--namespace is required")
	}
	if spec.Duration == "" {
		return fmt.Errorf("--duration is required")
	}
	if dryRunFlag := cmd.Flags().Lookup("dry-run"); dryRunFlag != nil && dryRunFlag.Changed && dryRunFlag.Value.String() == "true" {
		if ds, ok := s.(dryRunSubmitter); ok {
			return ds.DryRun(context.Background(), spec, cmd.OutOrStdout())
		}
		return NewPrintSubmitter(cmd.OutOrStdout()).Submit(context.Background(), spec)
	}
	return s.Submit(context.Background(), spec)
}
