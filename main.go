package main

import (
	"html/template"
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
    <script type="text/javascript">
        var ws = new WebSocket("ws://{{.}}/ws");
        ws.onopen = function() {
            console.log("Connected to WebSocket");
            ws.send("Hello Server");
        };
        ws.onmessage = function(evt) {
            var received_msg = evt.data;
            console.log("Message is received:", received_msg);
        };
        ws.onclose = function() { 
            console.log("Connection is closed..."); 
        };
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

		// Write message back to browser
		if err = conn.WriteMessage(msgType, msg); err != nil {
			return
		}
	}
}

func main() {
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
