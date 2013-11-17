package main

import (
	"fmt"
	"net/http"
	"os"
	"draringi/codejam2013/src/web"
)

func main() {
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	machine := new(web.MachineInterface)
	http.Handle("/upload", machine)
    dashboard := new(web.Dashboard)
    go dashboard.Init()
	http.Handle("/", dashboard)
    http.Handle("/data", dashboard.JSONAid)
	fmt.Println("listening...")
	err := http.ListenAndServe(os.Getenv("HOST")+":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}
