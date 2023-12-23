package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Define a template for the HTML page
var html = `
<!DOCTYPE html>
<html>
<head>
    <title>WebSocket Test</title>
</head>
<body>
    <h1>WebSocket Test</h1>
    <input type="text" id="messageInput" placeholder="Type a message...">
    <button onclick="sendMessage()">Send</button>
    <script type="text/javascript">
        var ws = new WebSocket("ws://{{.}}/ws");
        ws.onopen = function() {
            console.log("Connected to WebSocket");
        };
        ws.onmessage = function(evt) {
            var received_msg = evt.data;
            console.log("Message is received:", received_msg);
			document.body.innerHTML = evt.data; // Replace the entire body, or target a specific element
        };
        ws.onclose = function() { 
            console.log("Connection is closed..."); 
        };
        function sendMessage() {
            var message = document.getElementById('messageInput').value;
            ws.send(message);
            console.log("Message sent:", message);
        }
    </script>
</body>
</html>
`

// WebSocket Upgrader
var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {

		return
	}
	defer conn.Close()

	for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// Print the message to the console
		println(string(msg))
		if string(msg) == "katy" {
			print("Hey that's where I live!")
			// The URL you want to make a request to
			url := "http://wttr.in"

			// Make a GET request
			resp, err := http.Get(url)
			if err != nil {
				// Handle error
				fmt.Println(err)
				return
			}
			defer resp.Body.Close() // Make sure to close the body

			// Read the response body using io.ReadAll
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				// Handle error
				fmt.Println(err)
				return
			}

			// Print the response body
			fmt.Println(string(body))

		}

		// Write message back to browser
		if err = conn.WriteMessage(msgType, msg); err != nil {
			return
		}
	}
}

func main() {

	///Websockets begins here.

	r := gin.Default()

	// Serve HTML page at root
	r.GET("/", func(c *gin.Context) {
		t := template.New("webpage")
		t, _ = t.Parse(html)
		t.Execute(c.Writer, c.Request.Host)
	})

	// Handle WebSocket connections
	r.GET("/ws", func(c *gin.Context) {
		websocketHandler(c.Writer, c.Request)
	})

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
