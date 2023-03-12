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
    socket.send(input.value);
    input.value = "";
}