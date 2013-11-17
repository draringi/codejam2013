package web

import (
	"net/http"
	"draringi/codejam2013/src/forecasting"
	"draringi/codejam2013/src/data"
//	"html/template"
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
	self.channel = make(chan (*data.CSVData), 1)
	forecasting.PredictPulse(self.channel)
	go func () {
		for {
			tmp := <-self.channel
            if tmp != nil {
				self.data = tmp
			}
		}
	} ()
}

func (self *Dashboard) ServeHTTP (w http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(w, "Placeholder")
}
