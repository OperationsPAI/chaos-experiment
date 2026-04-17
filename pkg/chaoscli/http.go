package chaoscli

import "github.com/spf13/cobra"

func NewHTTPCmd(s Submitter, opts *commandOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "http",
		Short: "HTTP fault commands",
	}
	cmd.AddCommand(newHTTPAbortCmd(s, opts))
	return cmd
}

func newHTTPAbortCmd(s Submitter, opts *commandOptions) *cobra.Command {
	var app string
	var duration string
	var route string
	var method string
	var port int

	cmd := &cobra.Command{
		Use:   "abort",
		Short: "Abort HTTP requests",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return submitSpec(cmd, s, Spec{
				Type:      "HTTPRequestAbort",
				Namespace: opts.namespace,
				Target:    app,
				Duration:  duration,
				Params: map[string]any{
					"route":       route,
					"http_method": method,
					"port":        port,
				},
			})
		},
	}

	cmd.Flags().StringVar(&app, "app", "", "Application label")
	cmd.Flags().StringVar(&duration, "duration", "", "Fault duration, for example 2m")
	cmd.Flags().StringVar(&route, "route", "", "HTTP route pattern")
	cmd.Flags().StringVar(&method, "http-method", "GET", "HTTP method")
	cmd.Flags().IntVar(&port, "port", 80, "Service port")
	return cmd
}
