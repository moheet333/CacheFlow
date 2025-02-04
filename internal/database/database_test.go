package database

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Setenv("CACHEFLOW_REDIS_ADDRESS", "localhost")
	os.Setenv("CACHEFLOW_REDIS_PORT", "6379")
	os.Setenv("CACHEFLOW_REDIS_PASSWORD", "")
	os.Setenv("CACHEFLOW_REDIS_DATABASE", "0")
	code := m.Run()
	os.Exit(code)
}

func TestNew(t *testing.T) {
	srv := New()
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestHealth(t *testing.T) {
	srv := New()

	stats := srv.Health()

	if stats["redis_status"] != "up" {
		t.Fatalf("expected status to be up, but got %s", stats["redis_status"])
	}

	if _, ok := stats["redis_version"]; !ok {
		t.Fatalf("expected redis_version to be present, got %v", stats["redis_version"])
	}
}
