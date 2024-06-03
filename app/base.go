package app

import "log/slog"

// BaseApp implements app.App and defines the app structure.
type BaseApp struct {
	// configurable parameters
	dataDir string
	isDev   bool

	// internals
	logger *slog.Logger
}

// BaseAppConfig defines a BaseApp configuration options.
type BaseAppConfig struct {
	IsDev   bool
	DataDir string
}

var _ App = (*BaseApp)(nil)

// NewBaseApp creates a new BaseApp instance.
func NewBaseApp(config BaseAppConfig) *BaseApp {
	app := &BaseApp{
		isDev:   config.IsDev,
		dataDir: config.DataDir,
	}

	return app
}

// IsDev returns true if the app is running in development mode.
func (app *BaseApp) IsDev() bool {
	return app.isDev
}

// DataDir return the app data directory.
func (app *BaseApp) DataDir() string {
	return app.dataDir
}

// Logger returns the default app logger.
func (app *BaseApp) Logger() *slog.Logger {
	if app.logger == nil {
		app.logger = slog.Default()
	}

	return app.logger
}
