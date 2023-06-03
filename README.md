# httpmid

The httpmid package provides a function that allows you to create middleware for the standard net/http library.

## Installation

```shell
go get github.com/pkg-id/httpmid
```

## Example

```go
package main

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/pkg-id/httpmid"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/demo", demoHandler)

	// Creates a new handler with middleware.
	// The requestID becomes the outermost layer, followed by the logger, and
	// the innermost layer is the actual handler (mux) where the API is registered.
	handler := httpmid.Reduce(requestID, logger).Then(mux)

	err := http.ListenAndServe("localhost:8080", handler)
	if err != nil {
		log.Fatalf("listen failed: %v", err)
	}
}

func demoHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "Hello, World")
}

const headerRequestID = "X-Request-ID"

func requestID(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get(headerRequestID)
		if rid == "" {
			b := make([]byte, 32)
			_, _ = io.ReadFull(rand.Reader, b)
			rid = base64.URLEncoding.EncodeToString(b)
		}
		r.Header.Set(headerRequestID, rid)
		w.Header().Set(headerRequestID, rid)
		next.ServeHTTP(w, r)
	})
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		rid := r.Header.Get(headerRequestID)
		log.Printf("request started, request_id: %s, method: %s, path: %s", rid, r.Method, r.URL.Path)
		defer log.Printf("request completed, request_id: %s, latency: %s", rid, time.Since(started))

		next.ServeHTTP(w, r)
	})
}
```