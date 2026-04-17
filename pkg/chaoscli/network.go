package chaoscli

import "github.com/spf13/cobra"

func NewNetworkCmd(s Submitter, opts *commandOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network",
		Short: "Network fault commands",
	}

	cmd.AddCommand(
		newNetworkDelayCmd(s, opts),
		newNetworkLossCmd(s, opts),
	)

	return cmd
}

func newNetworkDelayCmd(s Submitter, opts *commandOptions) *cobra.Command {
	var app string
	var duration string
	var targetService string
	var latency int
	var correlation int
	var jitter int
	var direction string

	cmd := &cobra.Command{
		Use:   "delay",
		Short: "Inject network delay",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return submitSpec(cmd, s, Spec{
				Type:      "NetworkDelay",
				Namespace: opts.namespace,
				Target:    app,
				Duration:  duration,
				Params: map[string]any{
					"target_service": targetService,
					"latency":        latency,
					"correlation":    correlation,
					"jitter":         jitter,
					"direction":      direction,
				},
			})
		},
	}

	cmd.Flags().StringVar(&app, "app", "", "Source application label")
	cmd.Flags().StringVar(&duration, "duration", "", "Fault duration, for example 2m")
	cmd.Flags().StringVar(&targetService, "target-service", "", "Destination service label")
	cmd.Flags().IntVar(&latency, "latency", 100, "Delay latency in milliseconds")
	cmd.Flags().IntVar(&correlation, "correlation", 0, "Delay correlation percentage")
	cmd.Flags().IntVar(&jitter, "jitter", 0, "Delay jitter in milliseconds")
	cmd.Flags().StringVar(&direction, "direction", "to", "Traffic direction: to|from|both")
	return cmd
}

func newNetworkLossCmd(s Submitter, opts *commandOptions) *cobra.Command {
	var app string
	var duration string
	var targetService string
	var loss int
	var correlation int
	var direction string

	cmd := &cobra.Command{
		Use:   "loss",
		Short: "Inject network packet loss",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return submitSpec(cmd, s, Spec{
				Type:      "NetworkLoss",
				Namespace: opts.namespace,
				Target:    app,
				Duration:  duration,
				Params: map[string]any{
					"target_service": targetService,
					"loss":           loss,
					"correlation":    correlation,
					"direction":      direction,
				},
			})
		},
	}

	cmd.Flags().StringVar(&app, "app", "", "Source application label")
	cmd.Flags().StringVar(&duration, "duration", "", "Fault duration, for example 2m")
	cmd.Flags().StringVar(&targetService, "target-service", "", "Destination service label")
	cmd.Flags().IntVar(&loss, "loss", 10, "Packet loss percentage")
	cmd.Flags().IntVar(&correlation, "correlation", 0, "Loss correlation percentage")
	cmd.Flags().StringVar(&direction, "direction", "to", "Traffic direction: to|from|both")
	return cmd
}
