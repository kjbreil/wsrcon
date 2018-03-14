# wsrcon
Websocket Rcon implemetation for Rust Experimental

### notes
* Do not use right now, api is likely to change drastically
* There is a test file, it doesn't work though

### todo

* unmarshall chat json

### example

Create docker image from command line

```bash
docker run --rm -p 28016:28016/tcp -p 28015:28015/tcp -p 28015:28015/udp --name rust-server kjbreil/rust-server
```

go get

```bash
go get github.com/kjbreil/wsrcon
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

    rcon.AddGenericHandler(basicGenericHandler)

    rcon.Start()
}

func basicGenericHandler(msg string) {
    fmt.Printf("Generic Message: %s", msg)
}

```