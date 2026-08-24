package router

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	logmapubsub "github.com/xd-dash/logma-serverless/pubsub"
	logmarouter "github.com/xd-dash/logma-serverless/router"

	"github.com/dash-xd/pubsub-serverless-single-threaded/redisconn"
)

// requestedChannels returns the channel names a client asked to
// subscribe to via repeated ?channel= query parameters, same as
// logma-serverless's own /run and /events routes.
func requestedChannels(r *http.Request) []string {
	return r.URL.Query()["channel"]
}

// subscribeHandler backs POST /subscribe. It behaves exactly like
// logma-serverless's POST /run: claim a Runtime, subscribe it to the
// requested channels, block until it stops (control:shutdown or client
// disconnect), and report the outcome as JSON. The one difference from
// logma-serverless's own Runtime construction is the Redis client: it's
// resolved per request via redisconn (from the container's environment,
// or a caller-supplied header pair) rather than fixed once from the
// container's own environment at process start, so this deployment can
// serve subscribers against whichever Redis instance each request names.
func subscribeHandler(holder *logmapubsub.Holder[*logmarouter.Runtime]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client, err := redisconn.Client(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		rt, ok := holder.Claim()
		if !ok {
			http.Error(w, "runtime already claimed", http.StatusConflict)
			return
		}
		rt.Client = client
		rt.RecordInvocation(r, middleware.GetReqID(r.Context()))
		rt.Subscribe(requestedChannels(r))

		rt.Start(r.Context())

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "stopped"})
	}
}
