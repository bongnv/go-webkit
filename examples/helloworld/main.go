package main

import (
	"context"
	"log"

	"github.com/bongnv/go-webkit"
)

// Request is an example of Request DTO.
type Request struct {
	Name string
}

func main() {
	app := webkit.New()
	app.GET("/hello-world/:name", func(ctx context.Context, req webkit.Request) error {
		reqDto := &Request{}
		if err := req.Decode(reqDto); err != nil {
			return err
		}

		return req.Respond("Hello " + reqDto.Name)
	})

	log.Println(app.Run())
}
