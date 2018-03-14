# wsrcon
Websocket Rcon implemetation for Rust Experimental


WIP do not use in production, api is likely to change drastically

Tests do not work right now
## Example

Create docker image from command line

```bash
docker run --rm -p 28016:28016/tcp -p 28015:28015/tcp -p 28015:28015/udp --name rust-server kjbreil/rust-server
```

Sample go code

```go

package main

import (
    "fmt"

    "github.com/kjbreil/wsrcon"
)

func main() {
    // Connect to local docker
    ss := wsrcon.Settings{Host: "127.0.0.1", Port: 28016, Password: "docker"}

    rcon := wsrcon.Connect(&ss)

    rcon.AddChatHandler(basicGenericHandler)

    rcon.Start()
}

func basicGenericHandler(msg string) {
    fmt.Printf("Generic Message: %s", msg)
}

```