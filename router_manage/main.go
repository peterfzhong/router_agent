package main

import (
	"net/http"
	"log"
)

func main() {
	http.Handle("/css/", http.FileServer(http.Dir("template")))
	http.Handle("/js/", http.FileServer(http.Dir("template")))

	http.HandleFunc("/router_manage/", RouterManageHandler)

	log.Println("server start begin.........")
	http.ListenAndServe(":8888", nil)
	log.Println("server start end.........")
}