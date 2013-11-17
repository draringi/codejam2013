package main

import (
	"fmt"
	"net/http"
	"os"
	"draringi/codejam2013/src/web"
)

func main() {
    dashboard := new(web.Dashboard)
	machine := new(web.MachineInterface)
    go dashboard.Init()
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/upload", machine)
    dashboard.Lock.Lock()
    dashboard.Lock.Unlock()
	http.Handle("/", new(web.Static))
    http.Handle("/data", dashboard)
	fmt.Println("listening...")
	err := http.ListenAndServe(os.Getenv("HOST")+":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}
