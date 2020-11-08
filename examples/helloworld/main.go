package main

import (
	"context"
	"log"

	"github.com/bongnv/nanny"
)

// Request is an example of Request DTO.
type Request struct {
	Name string
}

func main() {
	app := nanny.Default()
	app.GET("/hello-world/:name", func(ctx context.Context, req nanny.Request) (interface{}, error) {
		reqDto := &Request{}
		if err := req.Decode(reqDto); err != nil {
			return nil, err
		}

		panic("I don't care")
		return "Hello " + reqDto.Name, nil
	})

	log.Println(app.Run())
}
