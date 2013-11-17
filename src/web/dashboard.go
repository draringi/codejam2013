package web

import (
	"net/http"
	"draringi/codejam2013/src/forecasting"
	"draringi/codejam2013/src/data"
    "encoding/json"
    "time"
)

type future struct {
    Records []record
}

type record struct {
    Date time.Time
    Power float64
}

type dashboardHelper struct {
    Data *data.CSVData
    Forcast *future
}

type Dashboard struct {
	channel chan (*data.CSVData)
	JSONAid *dashboardHelper
}

func (self *Dashboard) Init () {
	self.channel = make(chan (*data.CSVData), 1)
    JSONAid = new(dashboardHelper)
	forecasting.PredictPulse(self.channel)
	go func () {
		for {
			tmp := <-self.channel
            if tmp != nil {
				self.JSONAid.Data = tmp
			}
		}
	} ()
}

func (self *Dashboard) ServeHTTP (w http.ResponseWriter, request *http.Request) {
	http.ServeFile(w, request, "dashboard.html")
}

func (self *dashboardHelper) Build (Data *data.CSVData) {
    self.Data = Data
    self.Forcast = new(future)
    self.Forcast.Records = make([]record,len(Data.Data))
    for i :=0; i<len(Data.Data); i++ {
        self.Forcast.Records[i].Date = Data.Data[i].Time
        self.Forcast.Records[i].Power = Data.Data[i].Power
    }
}

func (self *dashboardHelper) jsonify (w http.ResponseWriter) {
    encoder := json.NewEncoder(w)
    encoder.Encode(self.Forcast)
}

func (self *dashboardHelper) ServeHTTP (w http.ResponseWriter, request *http.Request) {
    self.jsonify(w)
}
