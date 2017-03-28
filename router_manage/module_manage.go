package main

import (
	"net/http"
	"strings"
	"log"
)
func RouterManageHandler (w http.ResponseWriter, r *http.Request)  {
	pathInfo := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(pathInfo, "/")

	log.Println(pathInfo)
	log.Println(parts)
	var action = ""
	if len(parts) > 1 {
		action = strings.Title(parts[1])
	}


	if strings.EqualFold(action , "add_module") {
		AddRouterManageHandler(w, r)
		return
	}else if strings.EqualFold(action , "delete_module") {
		AddRouterManageHandler(w, r)
		return
	}else if strings.EqualFold(action , "update_module"){
		AddRouterManageHandler(w, r)
		return
	}else if strings.EqualFold(action , "query_module_list") {
		AddRouterManageHandler(w, r)
		return
	}

}

func AddRouterManageHandler (w http.ResponseWriter, r *http.Request) {
	log.Println("hello world")
	w.Write([] byte("hello world"))
}

func UpdateRouterManageHandler (w http.ResponseWriter, r *http.Request) {

}

func DelRouterManageHandler (w http.ResponseWriter, r *http.Request) {

}

func QueryRouterManageHandler (w http.ResponseWriter, r *http.Request) {

}