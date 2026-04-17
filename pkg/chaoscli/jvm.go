package chaoscli

import "github.com/spf13/cobra"

func NewJVMCmd(s Submitter, opts *commandOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jvm",
		Short: "JVM fault commands",
	}
	cmd.AddCommand(newJVMLatencyCmd(s, opts))
	return cmd
}

func newJVMLatencyCmd(s Submitter, opts *commandOptions) *cobra.Command {
	var app string
	var duration string
	var className string
	var methodName string
	var latencyDuration int

	cmd := &cobra.Command{
		Use:   "latency",
		Short: "Inject JVM latency",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return submitSpec(cmd, s, Spec{
				Type:      "JVMLatency",
				Namespace: opts.namespace,
				Target:    app,
				Duration:  duration,
				Params: map[string]any{
					"class":            className,
					"method":           methodName,
					"latency_duration": latencyDuration,
				},
			})
		},
	}

	cmd.Flags().StringVar(&app, "app", "", "Application label")
	cmd.Flags().StringVar(&duration, "duration", "", "Fault duration, for example 2m")
	cmd.Flags().StringVar(&className, "class", "", "Fully qualified JVM class")
	cmd.Flags().StringVar(&methodName, "method", "", "JVM method")
	cmd.Flags().IntVar(&latencyDuration, "latency-duration", 1000, "Latency in milliseconds")
	return cmd
}
