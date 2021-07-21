package main

import (
	"gowiki/app"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/view/", app.MakeHandler(app.ViewHandler))
	http.HandleFunc("/edit/", app.MakeHandler(app.EditHandler))
	http.HandleFunc("/save/", app.MakeHandler(app.SaveHandler))
	http.HandleFunc("/not_found", app.NotFoundHandler)
	http.HandleFunc("/", app.MainHandler)
	log.Println("Server hosted in http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
