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
	app := webkit.Default()
	app.GET("/hello-world/:name", func(ctx context.Context, req webkit.Request) (interface{}, error) {
		reqDto := &Request{}
		if err := req.Decode(reqDto); err != nil {
			return nil, err
		}

		return "Hello " + reqDto.Name, nil
	})

	log.Println(app.Run())
}
