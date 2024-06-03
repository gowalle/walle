package app

import (
	"os"
	"testing"
)

func TestNewBaseApp(t *testing.T) {
	const testDataDir = "./base_app_test_data_dir/"
	defer os.RemoveAll(testDataDir)

	app := NewBaseApp(BaseAppConfig{
		DataDir: testDataDir,
		IsDev:   true,
	})

	if app.dataDir != testDataDir {
		t.Fatalf("expected dataDir %q, got %q", testDataDir, app.dataDir)
	}

	if !app.isDev {
		t.Fatalf("expected isDev true, got %v", app.isDev)
	}
}
