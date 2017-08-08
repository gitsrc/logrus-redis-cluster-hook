package logredis

import (
	"os"
	"testing"
)

func TestNewHookFunc(t *testing.T) {
	config := HookConfig{
		Addrs:      []string{"localhost:6379"},
		ConnOption: SINGLE,
		Key:        "key",
		Format:     "format",
		App:        "appname",
		Password:   "password",
		DB:         1,
	}

	hook, err := NewHook(config)
	if hook != nil {
		t.Fatalf("TestNewHookFunc, expected no hook, got hook: %s", hook)
	}

	if err == nil {
		t.Fatalf("TestNewHookFunc, expected %q, got %s.", "unknown message format", err)
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
