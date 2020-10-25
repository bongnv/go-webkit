# go-webkit

[![Build](https://github.com/bongnv/go-webkit/workflows/Build/badge.svg)](https://github.com/bongnv/go-webkit/actions?query=workflow%3ABuild)
[![codecov](https://codecov.io/gh/bongnv/go-webkit/branch/main/graph/badge.svg?token=0SSLExlCNY)](https://codecov.io/gh/bongnv/go-webkit)
[![Go Report Card](https://goreportcard.com/badge/github.com/bongnv/go-webkit)](https://goreportcard.com/report/github.com/bongnv/go-webkit)

A webkit for Go with simple APIs to use. It solves common problems of a web server so engineers can focus on business logic.

## Features

- Graceful shutdown

## Quick Start
### Installation
```sh
go get github.com/bongnv/go-webkit
```

### Example

```go
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
```
