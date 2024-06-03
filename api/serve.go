package api

import (
	"net/http"

	"github.com/gowalle/walle/app"
)

// ServeConfig is the configuration for the serve command.
type ServeConfig struct {

	// HttpAddr is the TCP address to listen for the HTTP server (eg. `127.0.0.1:80`).
	HttpAddr string

	// HttpsAddr is the TCP address to listen for the HTTPS server (eg. `127.0.0.1:443`).
	HttpsAddr string

	// Optional domains list to use when issuing the TLS certificate.
	//
	// If not set, the host from the bound server address will be used.
	//
	// For convenience, for each "non-www" domain a "www" entry and
	// redirect will be automatically added.
	CertificateDomains []string

	// AllowedOrigins is an optional list of CORS origins (default to "*").
	AllowedOrigins []string
}

func Serve(a app.App, config ServeConfig) (*http.Server, error) {
	if len(config.AllowedOrigins) == 0 {
		config.AllowedOrigins = []string{"*"}
	}

	// TODO run migrates

	// TODO reload app settings

	// TODO init api

	server := &http.Server{}

	return server, server.ListenAndServe()
}
