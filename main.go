package main

import (
	"log"
	"net/http"
)

var addr = "localhost:8080"

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path == "/" {
		http.ServeFile(w, r,  "dist/index.html")
		return
	}
	http.ServeFile(w, r,  "dist" + r.URL.Path)
}

func main() {
	log.Println("Chat started")
	hub := newHub()
	go hub.run()
	//fs := http.FileServer(http.Dir("dist"))

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	http.HandleFunc("/", serveHome)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
