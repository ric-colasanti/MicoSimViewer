package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/gorilla/mux"
	"github.com/ric-colasanti/MicoSimViewer/microsimserver"
)

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		fmt.Println("Linux", url)
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Println(err)
	}

}

func main() {
	router := mux.NewRouter()

	// Serve static files
	s := http.StripPrefix("/js/", http.FileServer(http.Dir("./js/")))
	router.PathPrefix("/js/").Handler(s)

	router.HandleFunc("/echo", microsimserver.PushData)
	router.HandleFunc("/", microsimserver.Home)

	fmt.Println("Listening at port:5555")
	openbrowser("http://127.0.0.1:5555")
	http.ListenAndServe(":5555", router)
}
