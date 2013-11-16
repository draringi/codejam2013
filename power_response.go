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
	http.HandleFunc("/", hello)
	http.Handle("/upload", machine)
	http.Handle("/upload", dashboard)
	fmt.Println("listening...")
	err := http.ListenAndServe(os.Getenv("HOST")+":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func hello(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "hello, world")
}
