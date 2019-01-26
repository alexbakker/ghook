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
parse the payload. 

## github-hook-receiver

An example usecase for ghook is [cmd/github-hook-receiver](cmd/github-hook-receiver).
It's a simple HTTP service that receives hooks and handles them according to a
list of handlers defined in a config file. 

```
Usage of github-hook-receiver:
  -addr string
    	address to listen on (default "127.0.0.1:8080")
  -config string
    	the filename of the configuration file (default "config.json")
```

The configuration file should have the following format:

```json
{
    "secret": "d696f82431664d9ea93483789db0116c",
    "handlers": [
        {
            "repo": "alexbakker/site",
            "ref": "refs/heads/master",
            "command": "/usr/local/bin/update-blog"
        }
    ]
}
```
