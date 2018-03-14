package wsrcon

import (
	"fmt"
	"testing"
)

// doesn't actually test anything right now, just the structure for me to put logic in next
func TestRCON(t *testing.T) {
	ss := Settings{Host: "127.0.0.1", Port: 28016, Password: "docker"}

	rcon := Connect(&ss)

	rcon.AddGenericHandler(testBasicGenericHandler)

	rcon.Start()
}

func testBasicGenericHandler(msg string) {
	fmt.Printf("BASIC CHAT: %s", msg)
}
