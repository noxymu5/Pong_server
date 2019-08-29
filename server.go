package main

import (
	"flag"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"time"

	"./party"
	"./github.com/gorilla/websocket"
)

//default port and path values
var Port = "8080"
var Path = "./templates"
var tmpl *template.Template

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Host string
	}{r.Host}

	err := tmpl.Execute(w, data)

	if err != nil {
		fmt.Printf("Error! %s\n", err)
		return
	}
}

func connectionHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Printf("Error! %s\n", err)
		return
	}

	party.NewConnection(conn)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	h := http.NewServeMux()

	flag.StringVar(&Port, "port", Port, "Port on which the server will be launched.\n")
	flag.StringVar(&Path, "path", Path, "Path to the directory, containing your webpage, static \nfolder (for js, css, and other additional files) and assets folder (for PixelJS).\n")
	flag.Parse()

	h.HandleFunc("/", homeHandler)
	h.HandleFunc("/game", connectionHandler)
	tmpl = template.Must(template.ParseFiles(Path + "/main.html"))

	h.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(Path+"/static"))))
	h.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(Path+"/assets"))))

	err := http.ListenAndServe(":"+Port, h)

	if err != nil {
		fmt.Printf("Error! %s\n", err)
	}
}
