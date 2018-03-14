package wsrcon

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type received struct {
	Message    string
	Identifier int
	Type       string
	Stacktrace *string
}

// Chat is the message returned by rcon with the chat type
// UserId (lowercase last d) does not lint in golang so use json tag for just that one
// same data is present in single line format from Generic Type message
type Chat struct {
	Message  string
	UserID   int `json:"UserId"`
	Username string
	Color    string
	Time     int
}

// Settings is the connection settings for a rust server
type Settings struct {
	Host     string
	Port     int
	Password string
}

// RCON is a wrapped websocket connection
type RCON struct {
	conn           *websocket.Conn
	genericHandler *func(string)
	chatHandler    *func(Chat)
}

// AddGenericHandler adds a handler for generic messages and logs all such
// messages, further isolation of generic messages will be added later
func (r *RCON) AddGenericHandler(handlerFunction func(string)) {
	r.genericHandler = &handlerFunction
}

// AddChatHandler is passed the function that is run when a Chat type message is seen in rcon
func (r *RCON) AddChatHandler(handlerFunction func(Chat)) {
	r.chatHandler = &handlerFunction
}

// Start starts the watching for messages, you don't need to do this to send messages
func (r *RCON) Start() {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Using done to just trigger that the rcon connection is done
	done := make(chan struct{})

	go func() {
		// data is the json structure that comes out of rust rcon
		var data *received
		// defer the closing of the goroutine until done
		defer close(done)
		for {

			// buffering of messages logic should go here, create channel to hold messages
			// and annother channel to wait on handler to complete to send next message to handler

			// when there is a websocket message read that data
			err := r.conn.ReadJSON(&data)

			if err != nil {
				log.Println("read:", err)
				return
			}

			switch data.Type {
			case "Generic":
				if r.genericHandler != nil {
					gh := *r.genericHandler
					gh(data.Message)
				}
			case "Chat":
				var msg Chat
				if r.chatHandler != nil {
					ch := *r.chatHandler

					err := json.Unmarshal([]byte(data.Message), &msg)
					if err != nil {
						fmt.Printf("%s", err)
					}

					// fmt.Printf("Message: %s UserId: %d Username: %s Color: %s Time: %d", msg.Message, msg.UserID, msg.Username, msg.Color, msg.Time)
					ch(msg)
				}
			}

		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := r.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}

}

// Connect returns a connection for later working on
func Connect(s *Settings) (r RCON) {

	// combine port with address
	host := s.Host + ":" + strconv.Itoa(s.Port)

	// create URL scheme
	u := url.URL{Scheme: "ws", Host: host, Path: s.Password}
	log.Printf("connecting to %s", u.String())

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	r.conn = ws

	return
}

// Send a message to the server - this is the most basic of send function
// and can contain any type of precursor command
func (r *RCON) Send(msg string) {
	r.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}
