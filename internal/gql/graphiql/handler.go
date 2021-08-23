package graphiql

import (
	"log"
	"net/http"
)

// The endpoint is the url where you have your graphql api hosted
func New(endpoint string) http.HandlerFunc {
	t, err := preparingTemplate()
	if err != nil {
		log.Fatalln(err)
	}

	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_ = t.Execute(w, data(endpoint))
	}
}
