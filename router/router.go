// Package router builds pubsub-serverless-single-threaded's HTTP
// surface: a publish route and a subscribe route, each independently
// mountable so gospace-minimal can deploy them separately or compose
// them into one function (see NewRouter below, and the publishonly and
// subscribeonly subpackages). Both routes resolve their Redis
// connection per request via the redisconn package. The subscribe
// route reuses logma-serverless's own Runtime and Holder unchanged --
// it's the same claim/subscribe/block-until-stopped behavior
// logma-serverless's POST /run exposes, just mounted at /subscribe and
// without the paired SSE /events route, since this deployment never
// streams events back over HTTP.
package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	logmapubsub "github.com/xd-dash/logma-serverless/pubsub"
	logmarouter "github.com/xd-dash/logma-serverless/router"
)

// RegisterPublish mounts POST /publish on r.
func RegisterPublish(r chi.Router) {
	r.Post("/publish", publishHandler())
}

// RegisterSubscribe mounts POST /subscribe on r, backed by a
// logma-serverless Runtime holder scoped to this call -- one holder per
// mount, shared by every request the mounted route serves for the life
// of the container, the same way logma-serverless's own
// router.NewRouter wires its /run route.
func RegisterSubscribe(r chi.Router) {
	holder := logmapubsub.NewHolder(logmarouter.NewRuntime)
	r.Post("/subscribe", subscribeHandler(holder))
}

// NewRouter returns this deployment's default router: publish and
// subscribe mounted together on logma-serverless's shared chi
// middleware stack (see logmarouter.Build). To deploy publish and
// subscribe as separate functions instead, point gospace-minimal's
// router drop-in at the publishonly or subscribeonly subpackages, which
// each mount only one of the two routes on the same stack.
func NewRouter() http.Handler {
	return logmarouter.Build(func(r chi.Router) {
		RegisterPublish(r)
		RegisterSubscribe(r)
	})
}
