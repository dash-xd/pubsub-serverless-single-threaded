package function

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var mux = newMux()
func Main(w http.ResponseWriter, r *http.Request) {
	mux.ServeHTTP(w, r)
}
func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/subscribe", subscribe)
	return mux
}

func subscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		fmt.Printf("Received POST request with body: %s\n", body)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Request received successfully!")
}
