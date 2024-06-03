package walle

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/gowalle/walle/app"
	"github.com/spf13/cobra"
)

var Version = "(unknown)"

// appWrapper is a wrapper for app.App
type appWrapper struct {
	app.App
}

// Walle defines a Walle app launcher.
type Walle struct {
	*appWrapper

	devFlag     bool
	dataDirFlag string

	// RootCmd is the main command for walle.
	RootCmd *cobra.Command
}

// Config is the app initialization configuration.
type Config struct {
	DefaultDev     bool
	DefaultDataDir string
}

// New creates a new Walle instance with the default configuration.
func New() *Walle {
	_, isUsingGoRun := inspectRuntime()

	return NewWithConfig(Config{
		DefaultDev: isUsingGoRun,
	})
}

// NewWithConfig creates a new Walle instance with the specified configuration.
func NewWithConfig(config Config) *Walle {
	if config.DefaultDataDir == "" {
		baseDir, _ := inspectRuntime()
		config.DefaultDataDir = filepath.Join(baseDir, "data")
	}

	w := &Walle{
		RootCmd: &cobra.Command{
			Use:     filepath.Base(os.Args[0]),
			Short:   "Walle CLI",
			Version: Version,
			FParseErrWhitelist: cobra.FParseErrWhitelist{
				UnknownFlags: true,
			},
			// no need to provide the default cobra completion command
			CompletionOptions: cobra.CompletionOptions{
				DisableDefaultCmd: true,
			},
		},
		devFlag:     config.DefaultDev,
		dataDirFlag: config.DefaultDataDir,
	}

	// replace with a colored stderr writer
	w.RootCmd.SetErr(newErrWriter())

	// parse base flags
	// (errors are ignored, since the full flags parsing happens on Execute())
	w.eagerParseFlags(&config)

	// initialize the app instance
	w.appWrapper = &appWrapper{app.NewBaseApp(app.BaseAppConfig{
		IsDev:   w.devFlag,
		DataDir: w.dataDirFlag,
	})}

	// hide the default help command (allow only `--help` flag)
	w.RootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	return w
}

// eagerParseFlags parses the global app flags before calling pb.RootCmd.Execute().
// so we can have all Walle flags ready for use on initialization.
func (w *Walle) eagerParseFlags(config *Config) error {
	w.RootCmd.PersistentFlags().StringVar(
		&w.dataDirFlag,
		"dir",
		config.DefaultDataDir,
		"the Walle data directory",
	)

	w.RootCmd.PersistentFlags().BoolVar(
		&w.devFlag,
		"dev",
		config.DefaultDev,
		"enable dev mode, aka. printing logs and sql statements to the console",
	)

	return w.RootCmd.ParseFlags(os.Args[1:])
}

// Start starts the application, aka. registers the default system
// commands (serve, migrate, version) and executes pb.RootCmd.
func (w *Walle) Start() error {
	// register system commands
	// w.RootCmd.AddCommand(cmd.NewAdminCommand(pb))
	// w.RootCmd.AddCommand(cmd.NewServeCommand(pb, !pb.hideStartBanner))

	return w.Execute()
}

// Execute initializes the application (if not already) and executes
// the pb.RootCmd with graceful shutdown support.
//
// This method differs from pb.Start() by not registering the default
// system commands!
func (w *Walle) Execute() error {
	return nil
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

// newErrWriter returns a red colored stderr writter.
func newErrWriter() *coloredWriter {
	return &coloredWriter{
		w: os.Stderr,
		c: color.New(color.FgRed),
	}
}

// coloredWriter is a small wrapper struct to construct a [color.Color] writter.
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
