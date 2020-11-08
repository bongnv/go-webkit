# Usages
## Options
An `Option` customizes an `Application` and these are available `Option`:

### WithLogger
`WithLogger` allows to specify a custom implementation of the logger.
```go
  logger := log.New(os.Stderr, "", log.LstdFlags)
  app := nanny.New(nanny.WithLogger(logger))
```

### WithPProf

`WithPProf` starts another HTTP service in a different port to serve endpoints for [`pprof`](https://golang.org/pkg/net/http/pprof/). The option is included in the default app with port 8081.

```go
  app := nanny.New(WithPProf(":8081"))
```

## Route Options
A `RouteOption` customizes a route. It can be used to add middlewares like `Recovery()`.

For convenience, a `RouteOption` can be a `Option` for the `Application`. In this case, the `RouteOption` will be applied to all routes.

### WithRecovery
`WithRecovery` recovers from panics and returns error with 500 status code to clients.
```go
    app.GET("/hello-world", helloWorld, nanny.WithRecovery())
```

### WithDecoder
`WithDecoder` specifies a custom logic for decoding the request to a request DTO.
```go
   app.GET("/hello-world", helloWorld, nanny.WithDecoder(customDecoder))
```

### WithCORS
`WithCORS` enables the support for Cross-Origin Resource Sharing. Ref: https://developer.mozilla.org/en/docs/Web/HTTP/Access_control_CORS.
```go
  app := nanny.New(nanny.WithCORS(nanny.DefaultCORSConfig))
  // or
  app.GET("/hello-world", helloWorld, nanny.WithCORS(nanny.DefaultCORSConfig))
``` 

### WithTimeout
`WithTimeout` allows to specify the time limit for each route. 1 second timeout is included in the default app.
```go
  app.GET("/hello-world", helloWorld, nanny.WithTimeout(time.Second)
```

### WithErrorHandler
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

## Grouping routes

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

## Dependency injection

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
