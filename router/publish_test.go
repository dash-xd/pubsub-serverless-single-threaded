package router

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func newPublishRouter() *chi.Mux {
	r := chi.NewRouter()
	RegisterPublish(r)
	return r
}

func TestPublishRejectsMissingChannel(t *testing.T) {
	t.Setenv("REDIS_URI", "redis-env:6379")
	t.Setenv("REDISCLI_AUTH", "env-secret")

	r := newPublishRouter()

	req := httptest.NewRequest("POST", "/publish", strings.NewReader(`{"data":{"foo":"bar"}}`))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != 400 {
		t.Fatalf("expected 400 for missing channel, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestPublishRejectsInvalidJSON(t *testing.T) {
	t.Setenv("REDIS_URI", "redis-env:6379")
	t.Setenv("REDISCLI_AUTH", "env-secret")

	r := newPublishRouter()

	req := httptest.NewRequest("POST", "/publish", strings.NewReader(`not json`))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != 400 {
		t.Fatalf("expected 400 for invalid JSON, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestPublishRejectsMissingRedisAddress(t *testing.T) {
	t.Setenv("REDIS_URI", "")
	t.Setenv("REDISCLI_AUTH", "")

	r := newPublishRouter()

	req := httptest.NewRequest("POST", "/publish", strings.NewReader(`{"channel":"foo"}`))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != 400 {
		t.Fatalf("expected 400 for missing Redis address, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestPublishMethodNotAllowed(t *testing.T) {
	t.Setenv("REDIS_URI", "redis-env:6379")
	t.Setenv("REDISCLI_AUTH", "env-secret")

	r := newPublishRouter()

	req := httptest.NewRequest("GET", "/publish", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != 405 {
		t.Fatalf("expected 405 for GET /publish, got %d", rec.Code)
	}
}
