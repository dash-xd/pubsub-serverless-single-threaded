// Package publishonly is a gospace-minimal router drop-in that mounts
// only pubsub-serverless-single-threaded's publish route, for a
// deployment that never needs to subscribe.
package publishonly

import (
	"net/http"

	logmarouter "github.com/xd-dash/logma-serverless/router"

	"github.com/dash-xd/pubsub-serverless-single-threaded/router"
)

// NewRouter returns a router that mounts only POST /publish.
func NewRouter() http.Handler {
	return logmarouter.Build(router.RegisterPublish)
}
