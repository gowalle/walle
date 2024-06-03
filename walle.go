package walle

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// Version of Walle.
var Version = "(unknown)"

// App defines the interface for an application.
type App interface {
	// Logger returns the active app logger.
	Logger() *slog.Logger
}

type BaseApp struct {
	isDev   bool
	dataDir string

	// internals
	logger *slog.Logger

	// command
	RootCmd *cobra.Command
}

type Config struct {
	IsDev   bool
	DataDir string
}

var _ App = (*BaseApp)(nil)

// New creates a new BaseApp.
func New() *BaseApp {
	_, isUsingGoRun := inspectRuntime()
	return NewWithConfig(Config{IsDev: isUsingGoRun})
}

// NewWithConfig creates a new BaseApp with the given config.
func NewWithConfig(config Config) *BaseApp {
	if config.DataDir == "" {
		baseDir, _ := inspectRuntime()
		config.DataDir = filepath.Join(baseDir, "data")
	}
	app := &BaseApp{
		RootCmd: &cobra.Command{
			Use:     filepath.Base(os.Args[0]),
			Short:   "Walle CLI",
			Version: Version,
			FParseErrWhitelist: cobra.FParseErrWhitelist{
				UnknownFlags: false,
			},
			CompletionOptions: cobra.CompletionOptions{
				DisableDefaultCmd: true,
			},
		},
	}

	// replace with a colored stderr writer
	app.RootCmd.SetErr(newErrWriter())

	// parse base flags
	// (errors are ignored, since the full flags parsing happens on Execute())
	err := app.eagerParseFlags(&config)
	if err != nil {
		panic(err)
	}

	// hide the default help command (allow only `--help` flag)
	app.RootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	return app
}

// Start starts the application, aka. registers the default system
// commands (serve, migrate, version) and executes pb.RootCmd.
func (app *BaseApp) Start() error {
	// TODO register inner commands
	return app.Execute()
}

// Execute initializes the application (if not already) and executes
func (app *BaseApp) Execute() error {
	// TODO
	return nil
}

// Logger returns the active app logger.
func (app *BaseApp) Logger() *slog.Logger {
	if app.logger == nil {
		return slog.Default()
	}
	return app.logger
}

// eagerParseFlags parses the global app flags before calling pb.RootCmd.Execute().
// so we can have all PocketBase flags ready for use on initialization.
func (app *BaseApp) eagerParseFlags(config *Config) error {
	app.RootCmd.PersistentFlags().StringVar(
		&app.dataDir,
		"dir",
		config.DataDir,
		"the app data directory",
	)

	app.RootCmd.PersistentFlags().BoolVar(
		&app.isDev,
		"dev",
		config.IsDev,
		"enable dev mode, aka. printing logs and sql statements to the console",
	)

	return app.RootCmd.ParseFlags(os.Args[1:])
}

// inspectRuntime tries to find the base executable directory and how it was run.
func inspectRuntime() (baseDir string, withGoRun bool) {
	if strings.HasPrefix(os.Args[0], os.TempDir()) {
		// probably ran with go run
		withGoRun = true
		baseDir, _ = os.Getwd()
	} else {
		// probably ran with go build
		withGoRun = false
		baseDir = filepath.Dir(os.Args[0])
	}
	return
}

// newErrWriter returns a red colored stderr writer.
func newErrWriter() *coloredWriter {
	return &coloredWriter{
		w: os.Stderr,
		c: color.New(color.FgRed),
	}
}

// coloredWriter is a small wrapper struct to construct a [color.Color] writer.
type coloredWriter struct {
	w io.Writer
	c *color.Color
}

// Write writes the p bytes using the colored writer.
func (colored *coloredWriter) Write(p []byte) (n int, err error) {
	colored.c.SetWriter(colored.w)
	defer colored.c.UnsetWriter(colored.w)

	return colored.c.Print(string(p))
}
