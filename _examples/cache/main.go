package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tinyredglasses/workers2/cloudflare"
	"github.com/tinyredglasses/workers2/cloudflare/cache"
)

type responseWriter struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.StatusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	rw.Body = append(rw.Body, data...)
	return rw.ResponseWriter.Write(data)
}

func (rw *responseWriter) ToHTTPResponse() *http.Response {
	return &http.Response{
		StatusCode: rw.StatusCode,
		Header:     rw.Header(),
		Body:       io.NopCloser(bytes.NewReader(rw.Body)),
	}
}

func handler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	rw := responseWriter{ResponseWriter: w}
	c := cache.New()

	// Find cache
	res, _ := c.Match(req, nil)
	if res != nil {
		// Set the response status code
		rw.WriteHeader(res.StatusCode)
		// Set the response headers
		for key, values := range res.Header {
			for _, value := range values {
				rw.Header().Add(key, value)
			}
		}
		rw.Header().Add("X-Message", "cache from worker")
		// Set the response body
		io.Copy(rw.ResponseWriter, res.Body)
		return
	}

	// Responding
	text := fmt.Sprintf("time:%v\n", time.Now().UnixMilli())
	rw.Header().Set("Cache-Control", "max-age=15")
	rw.Write([]byte(text))

	// Create cache
	cloudflare.WaitUntil(ctx, func() {
		err := c.Put(req, rw.ToHTTPResponse())
		if err != nil {
			fmt.Println(err)
		}
	})
}

func main() {
	workers.Serve(http.HandlerFunc(handler))
}
