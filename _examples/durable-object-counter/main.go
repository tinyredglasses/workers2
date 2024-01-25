package main

import (
	"io"
	"net/http"

	"github.com/tinyredglasses/workers2/cloudflare"
)

func main() {
	workers.Serve(&MyHandler{})
}

type MyHandler struct{}

func (_ *MyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	COUNTER, err := cloudflare.NewDurableObjectNamespace(req.Context(), "COUNTER")
	if err != nil {
		panic(err)
	}

	id := COUNTER.IdFromName("A")
	obj, err := COUNTER.Get(id)
	if err != nil {
		panic(err)
	}

	res, err := obj.Fetch(req)
	if err != nil {
		panic(err)
	}

	count, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	w.Write([]byte("Durable object 'A' count: " + string(count)))
}
