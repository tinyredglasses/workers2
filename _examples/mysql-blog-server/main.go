package main

import (
	"net/http"

	"github.com/tinyredglasses/workers2/_examples/mysql-blog-server/app"
)

func main() {
	http.Handle("/articles", app.NewArticleHandler())
	workers.Serve(nil) // use http.DefaultServeMux
}
