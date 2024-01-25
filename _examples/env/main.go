package main

import (
	"fmt"
	"net/http"

	"github.com/tinyredglasses/workers2/cloudflare"
)

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "MY_ENV: %s", cloudflare.Getenv(req.Context(), "MY_ENV"))
	})
	workers.Serve(handler)
}
