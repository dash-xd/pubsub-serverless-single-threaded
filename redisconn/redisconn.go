// Package redisconn resolves which Redis instance a single HTTP request
// should talk to, so a publish or subscribe route never has to hardcode
// that logic itself.
package redisconn

import (
	"fmt"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
)

const (
	headerRedisURI     = "X-Redis-Uri"
	headerRedisCliAuth = "X-Rediscli-Auth"
)

// Config is one request's resolved Redis connection info.
type Config struct {
	Addr string
	Auth string
}

// Resolve determines a request's Redis connection settings:
//
//  1. REDIS_URI and REDISCLI_AUTH both set in the environment -> use
//     both, the same for every request this container serves.
//  2. Only REDIS_URI set in the environment -> use it, but take the
//     auth password from this request's X-Rediscli-Auth header, so one
//     fixed instance can still be reached with a per-caller password.
//  3. Neither set in the environment -> take both the address and the
//     auth password from this request's X-Redis-Uri and
//     X-Rediscli-Auth headers, so a deployment with no Redis instance
//     of its own can serve whichever instance each caller names.
//
// It returns an error if the resolved address ends up empty.
func Resolve(r *http.Request) (Config, error) {
	envURI := os.Getenv("REDIS_URI")
	envAuth := os.Getenv("REDISCLI_AUTH")

	cfg := Config{Addr: envURI, Auth: envAuth}
	switch {
	case envURI != "" && envAuth != "":
		// Both set: use the environment as-is.
	case envURI != "":
		cfg.Auth = r.Header.Get(headerRedisCliAuth)
	default:
		cfg.Addr = r.Header.Get(headerRedisURI)
		cfg.Auth = r.Header.Get(headerRedisCliAuth)
	}

	if cfg.Addr == "" {
		return Config{}, fmt.Errorf("no Redis address configured: set REDIS_URI, or send it via the %s header", headerRedisURI)
	}
	return cfg, nil
}

// Client builds a *redis.Client for the connection Resolve determines
// for r. It does not connect until the first command is issued.
func Client(r *http.Request) (*redis.Client, error) {
	cfg, err := Resolve(r)
	if err != nil {
		return nil, err
	}
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Auth,
		DB:       0,
	}), nil
}
