package router

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/dash-xd/pubsub-serverless-single-threaded/redisconn"
)

// publishRequest is the payload POSTed to /publish.
type publishRequest struct {
	Channel string          `json:"channel"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// relayedMessage is what actually gets published to Redis: shaped so
// that a logma-serverless Runtime subscribed to Channel decodes it as a
// PublishRequest with an empty Channel field, and so falls back to the
// Redis channel it arrived on -- the same channel this handler
// published to. See handlePublish in logma-serverless's router package.
type relayedMessage struct {
	Data json.RawMessage `json:"data,omitempty"`
}

// publishHandler backs POST /publish: it resolves which Redis instance
// to use for this request (see redisconn.Resolve) and publishes the
// request body's data onto the named channel.
func publishHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "error reading request body", http.StatusInternalServerError)
			return
		}

		var req publishRequest
		if len(body) > 0 {
			if err := json.Unmarshal(body, &req); err != nil {
				http.Error(w, "invalid JSON body", http.StatusBadRequest)
				return
			}
		}
		if req.Channel == "" {
			http.Error(w, "channel is required", http.StatusBadRequest)
			return
		}

		client, err := redisconn.Client(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer client.Close()

		message, err := json.Marshal(relayedMessage{Data: req.Data})
		if err != nil {
			http.Error(w, "failed to encode message", http.StatusInternalServerError)
			return
		}

		if err := client.Publish(r.Context(), req.Channel, message).Err(); err != nil {
			http.Error(w, "failed to publish: "+err.Error(), http.StatusBadGateway)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "published", "channel": req.Channel})
	}
}
