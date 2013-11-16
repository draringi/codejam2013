package main

import (
	"fmt"
	"net/http"
	"os"
	"./src/web"
)

func main() {
	machine := new(web.MachineInterface)
	http.HandleFunc("/", hello)
	http.Handle("/upload", machine)
	fmt.Println("listening...")
	err := http.ListenAndServe(os.Getenv("HOST")+":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func hello(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "hello, world")
}
