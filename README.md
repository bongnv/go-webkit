# gwf

[![Build](https://github.com/bongnv/gwf/workflows/Build/badge.svg)](https://github.com/bongnv/gwf/actions?query=workflow%3ABuild)
[![codecov](https://codecov.io/gh/bongnv/gwf/branch/main/graph/badge.svg?token=0SSLExlCNY)](https://codecov.io/gh/bongnv/gwf)
[![Go Report Card](https://goreportcard.com/badge/github.com/bongnv/gwf)](https://goreportcard.com/report/github.com/bongnv/gwf)

A webkit for Go with simple APIs to use. It solves common problems of a web server so engineers can focus on business logic.

## Features

- Graceful shutdown

## Quick Start
### Installation
```sh
go get github.com/bongnv/gwf
```

### Example

```go
package main

import (
	"context"
	"log"

	"github.com/bongnv/gwf"
)

func main() {
	app := gwf.Default()
	app.GET("/hello-world", func(ctx context.Context, req gwf.Request) (interface{}, error) {
		return "OK", nil
	})
	log.Println(app.Run())
}
```

## Usages
### Options
An `Option` customizes an `Application` and these are available `Option`:

#### WithLogger
`WithLogger` allows to specify a custom implementation of the logger.
```go
  logger := log.New(os.Stderr, "", log.LstdFlags)
  app := gwf.New(gwf.WithLogger(logger))
```

### Route Options
A `RouteOption` customizes a route. It can be used to add middlewares like `Recovery()`.

For convenience, a `RouteOption` can be a `Option` for the `Application`. In this case, the `RouteOption` will be applied to all routes.

#### WithRecovery
`WithRecovery` recovers from panics and returns error with 500 status code to clients.
```go
    app.GET("/hello-world", helloWorld, gwf.WithRecovery())
```

#### WithDecoder
`WithDecoder` specifies a custom logic for decoding the request to a request DTO.
```go
   app.GET("/hello-world", helloWorld, gwf.WithDecoder(customDecoder))
```

#### WithCORS
`WithCORS` enables the support for Cross-Origin Resource Sharing. Ref: https://developer.mozilla.org/en/docs/Web/HTTP/Access_control_CORS.
```go
  app := gwf.New(gwf.WithCORS(gwf.DefaultCORSConfig))
  // or
  app.GET("/hello-world", helloWorld, gwf.WithCORS(gwf.DefaultCORSConfig))
``` 

#### WithErrorHandler
`WithErrorHandler` allows to specify a custom `ErrorHandler` which converts an error into HTTP response.
```go
func yourCustomErrHandler(w http.ResponseWriter, errResp error) {
    w.WriteHeader(http.StatusInternalServerError)
    if err := encoder.Encode(w, errResp); err != nil {
        logger.Println("Error", err, "while encoding", errResp)
    }
}

func yourInitFunc(app *gwf.Application) {
    app.GET("/hello-world", helloWorld, gwf.WithErrorHandler(yourCustomErrHandler))
}
```
