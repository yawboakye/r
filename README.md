## Shall We Retry?

Some function calls, when they fail, are worth retrying. We had a couple of
these calls in our application so we built a small retry utility. It's small
because it still leaves some of the heavy-lifting to you. Most importantly,
interface conversion is left to you. Also it accepts only one function type, but
it's broad enough to cover your use case. All you need is a wrapper function
that knows how to call your real function.

It does not recover from `panic`s, neither does it abort remaining trials when
specific errors occur. If you need any of these, please feel free to submit a PR
to add them.

### Example

```go
package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/yawboakye/r"
	"github.com/yawboakye/r/backoff"
)

func main() {
	f := r.F{
		MaxRetries: 5,
		Backoff:    backoff.Exponential{time.Millisecond},
		Fn: func(...interface{}) (interface{}, error) {
			return retryable()
		},
	}

	if n, err := f.Run(); err != nil {
		fmt.Printf("retryable: no even number generated after max trials")
	} else {
		num := n.(int)
		fmt.Printf("%d generated after %d trials", num, f.Tried())
	}
}

func retryable() (int, error) {
	b := rand.Int()
	if b%2 != 0 {
		return b, errors.New("not an even number")
	}
	return b, nil
}
```
