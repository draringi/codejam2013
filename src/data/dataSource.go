package data

import (
	"time"
	"io"
)

const ISO = "2006-01-02T15:04-05:00"

type Record struct {
	Time time.Time
	Radiation float64
	Humidity float64
	Temperature float64
	Wind float64
	Power float64
	Empty bool
	Null bool
}

type CSVData struct {
	Labels []string
	Data []Record
}

func AddCSVToDB (file io.Reader) () {
	_, data := csvParse(file)
	creativeUpdate(data)
}

type CSVRequest struct {
	Request io.Reader
	Return chan (*CSVData)
}

type DataSource struct {
	CSVChan chan (*CSVRequest)
	PulseChan chan (chan *Record)
}

func CreateDataSource () (*DataSource) {
	data := new(DataSource)
	
	go func () {
		for {
			select {
			case csv := <-data.CSVChan:
				val := new(CSVData)
				val.Labels, val.Data = csvParse(csv.Request)
				csv.Return <-val
				
			case pulse := <-data.PulseChan:
				pulse <- new(Record)
			}
		}
	} ()
	return data
}
