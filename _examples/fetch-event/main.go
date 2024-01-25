package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/tinyredglasses/workers2/cloudflare"
	"github.com/tinyredglasses/workers2/cloudflare/fetch"
)

func handler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	cloudflare.PassThroughOnException(ctx)

	// logging after responding
	cloudflare.WaitUntil(ctx, func() {
		for i := 0; i < 5; i++ {
			time.Sleep(time.Second)
		}
		fmt.Println("5-second task completed")
	})

	// panic if x-error header has provided
	if req.Header.Get("x-error") != "" {
		panic("error")
	}

	// responds with origin server
	fc := fetch.NewClient()
	proxy := httputil.ReverseProxy{
		Transport: fc.HTTPClient(fetch.RedirectModeManual).Transport,
		Director: func(r *http.Request) {
			r.URL = req.URL
		},
	}

	proxy.ServeHTTP(w, req)
}

func main() {
	workers.Serve(http.HandlerFunc(handler))
}
