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
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(string(msgType))
			return
		}
		println(string(msg))
		if string(msg) == "katy" {
			// Fetch weather data
			url := "http://wttr.in"
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("Error fetching weather:", err)
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response:", err)
				return
			}

			// Send weather data to the client
			if err := conn.WriteMessage(websocket.TextMessage, body); err != nil {
				fmt.Println("Error sending message:", err)
				return
			}
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
