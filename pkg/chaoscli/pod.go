package chaoscli

import "github.com/spf13/cobra"

func NewPodCmd(s Submitter, opts *commandOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pod",
		Short: "Pod fault commands",
	}
	cmd.AddCommand(newPodFailureCmd(s, opts))
	return cmd
}

func newPodFailureCmd(s Submitter, opts *commandOptions) *cobra.Command {
	var app string
	var duration string

	cmd := &cobra.Command{
		Use:   "failure",
		Short: "Inject pod failure",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return submitSpec(cmd, s, Spec{
				Type:      "PodFailure",
				Namespace: opts.namespace,
				Target:    app,
				Duration:  duration,
			})
		},
	}

	cmd.Flags().StringVar(&app, "app", "", "Application label")
	cmd.Flags().StringVar(&duration, "duration", "", "Fault duration, for example 2m")
	return cmd
}
