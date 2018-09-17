# ghook [![godoc](https://godoc.org/github.com/alexbakker/ghook?status.svg)](https://godoc.org/github.com/alexbakker/ghook)

__ghook__ is a small Go package for receiving GitHub web hooks.

Example usage:

```go
package main

import (
    "fmt"
    "net/http"

    "github.com/alexbakker/ghook"
)

var (
    secret = []byte("d696f82431664d9ea93483789db0116c")
)

func main() {
    hook := ghook.New(secret, func(event *ghook.Event) error {
        fmt.Printf("received %s event!\n", event.Name)
        return nil
    })

    panic(http.ListenAndServe("127.0.0.1:8080", hook))
}
```

You can then use the [go-github](https://github.com/google/go-github) package to
parse the payload. See [cmd/github-hook-receiver](cmd/github-hook-receiver) for
a more elaborate example.
