package app

import "log/slog"

// App defines the interface for the application.
type App interface {

	// DataDir returns the app data directory path.
	DataDir() string

	// IsDev returns whether the app is in dev mode.
	IsDev() bool

	// Logger retruns the active app logger.
	Logger() *slog.Logger
}
