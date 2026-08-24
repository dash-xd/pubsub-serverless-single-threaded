// Package subscribeonly is a gospace-minimal router drop-in that mounts
// only pubsub-serverless-single-threaded's subscribe route, for a
// deployment that never needs to publish.
package subscribeonly

import (
	"net/http"

	logmarouter "github.com/xd-dash/logma-serverless/router"

	"github.com/dash-xd/pubsub-serverless-single-threaded/router"
)

// NewRouter returns a router that mounts only POST /subscribe.
func NewRouter() http.Handler {
	return logmarouter.Build(router.RegisterSubscribe)
}
