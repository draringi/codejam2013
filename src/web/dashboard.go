package web
e

import (
    "io"
	"net/http"
	"draringi/codejam2013/src/forecasting"
	"draringi/codejam2013/src/data"
    "encoding/json"
    "time"
    "sync"
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
}

type Dashboard struct {
	channel chan ([]data.Record)
	JSONAid dashboardHelper
    Forcast *future
    Data []data.Record
    Lock sync.Mutex
}

func (self *Dashboard) Init () {
	self.Lock.Lock()
	self.channel = make(chan ([]data.Record), 1)
	self.Data = nil
	self.Forcast = nil
	self.Lock.Unlock()
	go forecasting.PredictPulse(self.channel)
	for {
		tmp := <-self.channel
		if tmp != nil {
			self.Data = tmp
			self.Build()
		}
	}
}

type Static struct{
}

func (self *Static) ServeHTTP (w http.ResponseWriter, request *http.Request) {
	http.ServeFile(w, request, "dashboard.html")
}

func (self *Dashboard) Build () {
    Data := self.Data
    self.Forcast = new(future)
    self.Forcast.Records = make([]record,len(Data))
    for i :=0; i<len(Data); i++ {
        self.Forcast.Records[i].Date = Data[i].Time.Format(time.ANSIC)
        self.Forcast.Records[i].Power = Data[i].Power
    }
}

func (self *Dashboard) jsonify (w io.Writer) error {
    encoder := json.NewEncoder(w)
    if self.Data != nil {
        encoder.Encode(self.Forcast)
        return nil
    } else {
        return &dataError{"Error: Could not load data", time.Now()}
    }
}

func (self *Dashboard) ServeHTTP (w http.ResponseWriter, request *http.Request) {
    err := self.jsonify(w)
    if err != nil {
        http.Error(w,err.Error(), 404)
    }
}
