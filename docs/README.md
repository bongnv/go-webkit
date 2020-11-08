# Introduction

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
