package main

import (
	"fmt"
	"net/http"
	"os"
	"draringi/codejam2013/src/web"
)

func main() {
	machine := new(web.MachineInterface)
	dashboard := new(web.Dashboard)
    dashboard.Init()
	http.Handle("/upload", machine)
	http.Handle("/", dashboard)
    http.Handle("/data", dashboard.JSONAid)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("listening...")
	err := http.ListenAndServe(os.Getenv("HOST")+":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}
