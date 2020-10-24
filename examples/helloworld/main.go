package main

import (
	"context"
	"log"

	"github.com/bongnv/go-webkit"
)

func main() {
	app := webkit.New()
	app.GET("/hello-world", func(ctx context.Context, req webkit.Request) error {
		return req.Response("OK")
	})
	log.Println(app.Run())
}
