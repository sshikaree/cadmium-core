package webinterface

import (
	"html/template"
	"log"
	"net/http"
)

var (
	tmpls = template.Must(template.ParseGlob("./common/webinterface/templates/*"))
)

func getIndex(w http.ResponseWriter, r *http.Request) {
	tmpls.ExecuteTemplate(w, "index.html", nil)
	return
}

func StartWebInterface(port string) {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./common/webinterface/assets"))))
	http.HandleFunc("/", getIndex)

	log.Println("Starting web server on port:", port)
	http.ListenAndServe(":"+port, nil)
}
