package cmd

import (
	"github.com/gowalle/walle/app"
	"github.com/spf13/cobra"
)

// NewServeCommand creates and returns new command responsible for
// starting the default Walle web server.
func NewServeCommand(app app.App) *cobra.Command {
	var (
		allowedOrigins []string
		httpAddr       string
		httpsAddr      string
	)

	command := &cobra.Command{
		Use:          "serve [domain(s)]",
		Args:         cobra.ArbitraryArgs,
		Short:        "Starts the web server (default to 127.0.0.1:8090 if no domain is specified)",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				if httpAddr == "" {
					httpAddr = "127.0.0.1:80"
				}
				if httpsAddr == "" {
					httpsAddr = "127.0.0.1:443"
				}
			} else {
				if httpAddr == "" {
					httpAddr = "127.0.0.1:8090"
				}
			}
			// TODO: add support for gRPC/HTTP
			return nil
		},
	}

	command.PersistentFlags().StringSliceVar(
		&allowedOrigins,
		"origins",
		[]string{"*"},
		"CORS allowed domain origins list",
	)

	command.PersistentFlags().StringVar(
		&httpAddr,
		"http",
		"",
		"TCP address to listen for the HTTP server\n(if domain args are specified - default to 0.0.0.0:80, otherwise - default to 127.0.0.1:8090)",
	)

	command.PersistentFlags().StringVar(
		&httpsAddr,
		"https",
		"",
		"TCP address to listen for the HTTPS server\n(if domain args are specified - default to 0.0.0.0:443, otherwise - default to empty string, aka. no TLS)\nThe incoming HTTP traffic also will be auto redirected to the HTTPS version",
	)

	return command
}
