package web

import (
	"net/http"
	"draringi/codejam2013/src/forecasting"
	"draringi/codejam2013/src/data"
	"html/template"
	"fmt"
)

type dashboardHelper struct {

}

type Dashboard struct {
	channel chan (*data.CSVData)
	JSONAid dashboardHelper
	data *data.CSVData
}

func (self *Dashboard) Init () {
	self.channel = forecasting.PredictPulse()
	go func () {
		for {
			if tmp := <-self.channel {
				self.data = tmp
			}
		}
	} ()
}

func (self *Dashboard) ServeHTTP (w http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(w, "Placeholder")
}
