package web

import (
    "io"
	"net/http"
	"draringi/codejam2013/src/forecasting"
	"draringi/codejam2013/src/data"
    "encoding/json"
    "time"
)

type dataError struct {
    What string
    When time.Time
}

func (self *dataError) Error() string {
    return "["+self.When.Format(data.ISO)+"] " + self.What
}

type future struct {
    Records []record
}

type record struct {
    Date string
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
    self.JSONAid = new(dashboardHelper)
    self.JSONAid.Data = nil
    self.JSONAid.Forcast = nil
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
        self.Forcast.Records[i].Date = Data.Data[i].Time.Format(data.ISO)
        self.Forcast.Records[i].Power = Data.Data[i].Power
    }
}

func (self *dashboardHelper) jsonify (w io.Writer) error {
    encoder := json.NewEncoder(w)
    if self.Data != nil {
        encoder.Encode(self.Forcast)
        return nil
    } else {
        return &dataError{"Error: Could not load data", time.Now()}
    }
}

func (self *dashboardHelper) ServeHTTP (w http.ResponseWriter, request *http.Request) {
    err := self.jsonify(w)
    if err != nil {
        http.Error(w,err.Error(), 404)
    }
}
