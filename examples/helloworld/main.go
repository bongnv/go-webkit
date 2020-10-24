package main

import (
	"context"
	"log"

	"github.com/bongnv/go-webkit"
)

func main() {
	app := webkit.New()
	app.GET("/hello-world", func(ctx context.Context, req webkit.Request) error {
		return nil
	})
	log.Println(app.Run())
}
