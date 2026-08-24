package redisconn

import (
	"net/http/httptest"
	"testing"
)

func TestResolve(t *testing.T) {
	t.Run("both env vars set uses env for both, ignoring headers", func(t *testing.T) {
		t.Setenv("REDIS_URI", "redis-env:6379")
		t.Setenv("REDISCLI_AUTH", "env-secret")

		r := httptest.NewRequest("POST", "/publish", nil)
		r.Header.Set("X-Redis-Uri", "redis-header:6379")
		r.Header.Set("X-Rediscli-Auth", "header-secret")

		cfg, err := Resolve(r)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Addr != "redis-env:6379" || cfg.Auth != "env-secret" {
			t.Fatalf("expected env values, got %+v", cfg)
		}
	})

	t.Run("only REDIS_URI set uses env address and header auth", func(t *testing.T) {
		t.Setenv("REDIS_URI", "redis-env:6379")
		t.Setenv("REDISCLI_AUTH", "")

		r := httptest.NewRequest("POST", "/publish", nil)
		r.Header.Set("X-Rediscli-Auth", "header-secret")

		cfg, err := Resolve(r)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Addr != "redis-env:6379" || cfg.Auth != "header-secret" {
			t.Fatalf("expected env address + header auth, got %+v", cfg)
		}
	})

	t.Run("neither env var set uses both headers", func(t *testing.T) {
		t.Setenv("REDIS_URI", "")
		t.Setenv("REDISCLI_AUTH", "")

		r := httptest.NewRequest("POST", "/publish", nil)
		r.Header.Set("X-Redis-Uri", "redis-header:6379")
		r.Header.Set("X-Rediscli-Auth", "header-secret")

		cfg, err := Resolve(r)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Addr != "redis-header:6379" || cfg.Auth != "header-secret" {
			t.Fatalf("expected header values, got %+v", cfg)
		}
	})

	t.Run("only REDISCLI_AUTH env set still falls back to headers for both", func(t *testing.T) {
		t.Setenv("REDIS_URI", "")
		t.Setenv("REDISCLI_AUTH", "env-secret")

		r := httptest.NewRequest("POST", "/publish", nil)
		r.Header.Set("X-Redis-Uri", "redis-header:6379")
		r.Header.Set("X-Rediscli-Auth", "header-secret")

		cfg, err := Resolve(r)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Addr != "redis-header:6379" || cfg.Auth != "header-secret" {
			t.Fatalf("expected header values, got %+v", cfg)
		}
	})

	t.Run("no address anywhere is an error", func(t *testing.T) {
		t.Setenv("REDIS_URI", "")
		t.Setenv("REDISCLI_AUTH", "")

		r := httptest.NewRequest("POST", "/publish", nil)

		if _, err := Resolve(r); err == nil {
			t.Fatal("expected an error when no Redis address is configured")
		}
	})
}

func TestClientUsesResolvedAddr(t *testing.T) {
	t.Setenv("REDIS_URI", "redis-env:6379")
	t.Setenv("REDISCLI_AUTH", "env-secret")

	r := httptest.NewRequest("POST", "/publish", nil)

	client, err := Client(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer client.Close()

	if got := client.Options().Addr; got != "redis-env:6379" {
		t.Fatalf("expected addr redis-env:6379, got %q", got)
	}
}
