# nanny

[![Build](https://github.com/bongnv/nanny/workflows/Build/badge.svg)](https://github.com/bongnv/nanny/actions?query=workflow%3ABuild)
[![codecov](https://codecov.io/gh/bongnv/nanny/branch/main/graph/badge.svg?token=0SSLExlCNY)](https://codecov.io/gh/bongnv/nanny)
[![Go Report Card](https://goreportcard.com/badge/github.com/bongnv/nanny)](https://goreportcard.com/report/github.com/bongnv/nanny)

A web framework for Go with simple APIs to use. It solves common problems of a web server so engineers can focus on business logic.

## Features

- Graceful shutdown
- Panic recovery
- CORS
- Gzip compression
- Dependency injection

## Quick Start
### Installation
Make sure Go (**version 1.13+ is required**) is installed.
```sh
go get github.com/bongnv/nanny
```

### Example

```go
package main

import (
	"context"
	"log"

	"github.com/bongnv/nanny"
)

func main() {
    app := nanny.Default()
    app.GET("/hello-world", func(ctx context.Context, req nanny.Request) (interface{}, error) {
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
  app := nanny.New(nanny.WithLogger(logger))
```

### Route Options
A `RouteOption` customizes a route. It can be used to add middlewares like `Recovery()`.

For convenience, a `RouteOption` can be a `Option` for the `Application`. In this case, the `RouteOption` will be applied to all routes.

#### WithRecovery
`WithRecovery` recovers from panics and returns error with 500 status code to clients.
```go
    app.GET("/hello-world", helloWorld, nanny.WithRecovery())
```

#### WithDecoder
`WithDecoder` specifies a custom logic for decoding the request to a request DTO.
```go
   app.GET("/hello-world", helloWorld, nanny.WithDecoder(customDecoder))
```

#### WithCORS
`WithCORS` enables the support for Cross-Origin Resource Sharing. Ref: https://developer.mozilla.org/en/docs/Web/HTTP/Access_control_CORS.
```go
  app := nanny.New(nanny.WithCORS(nanny.DefaultCORSConfig))
  // or
  app.GET("/hello-world", helloWorld, nanny.WithCORS(nanny.DefaultCORSConfig))
``` 

#### WithTimeout
`WithTimeout` allows to specify the time limit for each route. 1 second timeout is included in the default app.
```go
  app.GET("/hello-world", helloWorld, nanny.WithTimeout(time.Second)
```

#### WithErrorHandler
`WithErrorHandler` allows to specify a custom `ErrorHandler` which converts an error into HTTP response.
```go
func yourCustomErrHandler(w http.ResponseWriter, errResp error) error {
    w.WriteHeader(http.StatusInternalServerError)
    if err := encoder.Encode(w, errResp); err != nil {
        logger.Println("Error", err, "while encoding", errResp)
        return err
    }

    return nil
}

func yourInitFunc(app *nanny.Application) {
    app.GET("/hello-world", helloWorld, nanny.WithErrorHandler(yourCustomErrHandler))
}
```

### Grouping routes

`nanny` supports grouping routes which share the same prefix or options for better readability.
```go
func main() {
    app := nanny.Default()

    // v1 group
    v1 := app.Group("/v1", WithV1Option())
    v1.GET("/hello-world", helloWorldV1)

    // v2 group
    v2 := app.Group("/v2", WithV2Option())
    v2.GET("/hello-world", helloWorldV2)

    log.Println(app.Run())
}
```

### Dependency injection

`nanny` makes dependency injection much easier by `Register`.
```go
// Service declares dependencies via `inject` tag.
type Service struct {
    DB *db.DB `inject:"db"` 
    Logger *Logger `inject:"logger"`
}

func main() {
    // Logger is provided with default implementation.
    app := nanny.Default()
    // Register DB component.
    app.Register("db", newDB())

    s := &Service{}
    // Logger and DB will be injected when registering service.
    app.Register("service", s)
}
```
