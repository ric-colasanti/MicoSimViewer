package microsimserver

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var conn *websocket.Conn
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Output(msgType int, tag string) {
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		fmt.Println("out")
		if err := conn.WriteMessage(msgType, []byte(fmt.Sprintf("%s  %d", tag, rand.Intn(100)))); err != nil {
			return
		}
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("home")
	// http.ServeFile(w, r, "index.html")
	fmt.Fprintf(w, `
	<!doctype html>
	<html lang="en">

	<head>
	  <title>MicosimViewer</title>
	  <!-- Required meta tags -->
	  <meta charset="utf-8">
	  <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

	  <!-- Bootstrap CSS -->
	  <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css"
		integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
	</head>

	<body>
	  <div class="container">
		<h3>MicrosimViewer static</h3>
		<button onclick="send()">Clickme</button>
		<pre id="output"></pre>

	  </div>
	  <!-- Optional JavaScript -->
	  <!-- jQuery first, then Popper.js, then Bootstrap JS -->
	  <script src="https://code.jquery.com/jquery-3.3.1.slim.min.js"
		integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo"
		crossorigin="anonymous"></script>
	  <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js"
		integrity="sha384-UO2eT0CpHqdSJQ6hJty5KVphtPhzWj9WO1clHTMGa3JDZwrnQq4sF86dIHNDz0W1"
		crossorigin="anonymous"></script>
	  <script src="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js"
		integrity="sha384-JjSmVgyd0p3pXB1rRibZUAYoIIy6OrQ6VrjIEaFf/nJGzIxFDsf4x0xIM+B07jRM"
		crossorigin="anonymous"></script>
	  <script>
	  console.log("MicrosimServer");
	  var input = document.getElementById("input");
	  var output = document.getElementById("output");
	  var socket = new WebSocket("ws://127.0.0.1:5555/echo");

	  socket.onopen = function () {
		  output.innerHTML += "Status: Connected\n";
	  };

	  socket.onmessage = function (e) {
		  output.innerHTML += "Server: " + e.data + "\n";
	  };

	  function send() {
		  socket.send("test");
		  input.value = "";
	  }
	  </script>
	</body>

	</html>
	 `)
}

func PushData(w http.ResponseWriter, r *http.Request) {
	conn, _ = upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
	fmt.Println("here")
	closeHandler := conn.CloseHandler()
	conn.SetCloseHandler(func(code int, text string) error {
		fmt.Println("closing")
		err := closeHandler(code, text)
		// ... or here.
		os.Exit(3)
		return err
	})

	for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		fmt.Println(msgType, err)

		// Print the message to the console
		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

		// Write message back to browser
		go Output(msgType, string(msg))
	}
}
