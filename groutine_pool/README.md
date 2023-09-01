# goroutine_pool

goroutine_pool is simple goroutine pool refer
to  [foundation](https://github.com/go-chassis/foundation/tree/master/gopool)

## Usage

### Sample code:

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/sunpe/gobox/goroutine_pool"
)

func main() {
	p := NewPool()
	p.Execute(func(ctx context.Context) {
		fmt.Println("goroutine poll demo")
	})
	p.Close(true)
}
```
