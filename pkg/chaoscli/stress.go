package chaoscli

import "github.com/spf13/cobra"

func NewStressCmd(s Submitter, opts *commandOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stress",
		Short: "Stress fault commands",
	}
	cmd.AddCommand(newStressCPUCmd(s, opts))
	return cmd
}

func newStressCPUCmd(s Submitter, opts *commandOptions) *cobra.Command {
	var app string
	var duration string
	var container string
	var cpuLoad int
	var cpuWorker int

	cmd := &cobra.Command{
		Use:   "cpu",
		Short: "Inject CPU stress",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return submitSpec(cmd, s, Spec{
				Type:      "CPUStress",
				Namespace: opts.namespace,
				Target:    app,
				Duration:  duration,
				Params: map[string]any{
					"container":  container,
					"cpu_load":   cpuLoad,
					"cpu_worker": cpuWorker,
				},
			})
		},
	}

	cmd.Flags().StringVar(&app, "app", "", "Application label")
	cmd.Flags().StringVar(&duration, "duration", "", "Fault duration, for example 2m")
	cmd.Flags().StringVar(&container, "container", "", "Container name override")
	cmd.Flags().IntVar(&cpuLoad, "cpu-load", 80, "CPU load percentage")
	cmd.Flags().IntVar(&cpuWorker, "cpu-worker", 1, "CPU worker count")
	return cmd
}
