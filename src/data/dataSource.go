package data

import (
	"time"
	"io"
)

const ISO = "2006-01-02T15:04Z05:00"

type Record struct {
	Time time.Time
	Radiation float64
	Humidity float64
	Temperature float64
	Wind float64
	Power float64
	empty bool
	null bool
}

type CVSData struct {
	Labels []string
	Data []Record
}

type CSVRequest struct {
	Request io.Reader
	Return chan (CVSData)
}

type DataSource struct {
	CSVChan chan (*CSVRequest)
	PulseChan chan (chan *Record)
}

func CreateDataSource () (DataSource) {
	var data DataSource
	
	go func () {
		for {
			select {
			case cvs := <-data.CSVChan:
				var val CVSData
				val.Labels, val.Data = csvParse(cvs.Request)
				cvs.Return <-val
				
			case pulse := <-data.PulseChan:
				pulse <- new(Record)
			}
		}
	} ()
	return data
}
