package router

import (
	"net/http/httptest"
	"testing"
)

func TestNewRouterMountsBothRoutes(t *testing.T) {
	t.Setenv("REDIS_URI", "")
	t.Setenv("REDISCLI_AUTH", "")

	r := NewRouter()

	for _, path := range []string{"/publish", "/subscribe"} {
		req := httptest.NewRequest("POST", path, nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code == 404 {
			t.Fatalf("expected %s to be mounted, got 404", path)
		}
	}
}
