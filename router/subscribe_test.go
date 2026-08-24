package router

import (
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func newSubscribeRouter() *chi.Mux {
	r := chi.NewRouter()
	RegisterSubscribe(r)
	return r
}

func TestRequestedChannelsReadsChannelQueryParam(t *testing.T) {
	t.Run("no query params yields no channels", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/subscribe", nil)
		if got := requestedChannels(r); len(got) != 0 {
			t.Fatalf("expected no channels, got %v", got)
		}
	})

	t.Run("repeated channel params are all returned in order", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/subscribe?channel=a&channel=b", nil)
		got := requestedChannels(r)
		if len(got) != 2 || got[0] != "a" || got[1] != "b" {
			t.Fatalf("expected [a b], got %v", got)
		}
	})
}

func TestSubscribeRejectsMissingRedisAddress(t *testing.T) {
	t.Setenv("REDIS_URI", "")
	t.Setenv("REDISCLI_AUTH", "")

	r := newSubscribeRouter()

	req := httptest.NewRequest("POST", "/subscribe", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != 400 {
		t.Fatalf("expected 400 for missing Redis address, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestSubscribeMethodNotAllowed(t *testing.T) {
	t.Setenv("REDIS_URI", "redis-env:6379")
	t.Setenv("REDISCLI_AUTH", "env-secret")

	r := newSubscribeRouter()

	req := httptest.NewRequest("GET", "/subscribe", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != 405 {
		t.Fatalf("expected 405 for GET /subscribe, got %d", rec.Code)
	}
}
