package web

import (
	"net/http"
	"../forecasting"
	"../data"
	"encoding/csv"
	"fmt"
//	"time"
	"strconv"
)


type MachineInterface struct {
	records *data.CSVData
	parser *data.DataSource
}

func recordToString(record data.Record) []string{
	stringRecord := make([]string, 6)
	stringRecord[0] = record.Time.Format(data.ISO)
	stringRecord[1] = strconv.FormatFloat(record.Radiation,'f', -1, 64)
	stringRecord[2] = strconv.FormatFloat(record.Humidity,'f', -1, 64)
	stringRecord[3] = strconv.FormatFloat(record.Temperature,'f', -1, 64)
	stringRecord[4] = strconv.FormatFloat(record.Wind,'f', -1, 64)
	stringRecord[5] = strconv.FormatFloat(record.Power,'f', -1, 64)
	return stringRecord
}

func (self *MachineInterface) ServeHTTP (w http.ResponseWriter, request *http.Request) {
	upload, uploadHeader, err := request.FormFile("file")
	if err != nil {
		fmt.Fprint(w, "Error: a file is needed")
		return
	}
	out := csv.NewWriter(w)
	if self.parser == nil {
		self.parser = data.CreateDataSource()
	}
	self.records = forecasting.PredictCSV(upload, parser.CSVChan)
	err = out.Write(self.records.Labels)
	for i := 0; i<len(self.records.Data); i++ {
		out.Write(recordToString(self.records.Data[i]))
	}
	out.Flush()
	return
}
	
