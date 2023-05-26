# logger

logger is a simple wrapper of [slog](https://pkg.go.dev/golang.org/x/exp/slog).

## Usage

### Sample code:

```go
package main

import "github.com/sunpe/gobox/logger"

func main() {
    logger.Info("hello world")
}
```

## Init

```go
package main

import "github.com/sunpe/gobox/logger"

func main() {
    logger.Init(loggger.WithLevel(logger.LevelDebug))
	logger.Info("hello world")
}
```
