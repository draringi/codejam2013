package web

import (
	"net/http"
	"draringi/codejam2013/src/forecasting"
	"draringi/codejam2013/src/data"
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
	stringRecord := make([]string, 2)
	stringRecord[0] = record.Time.Format(data.ISO)
	stringRecord[1] = strconv.FormatFloat(record.Power,'f', -1, 64)
	return stringRecord
}

func (self *MachineInterface) ServeHTTP (w http.ResponseWriter, request *http.Request) {
	if self.parser == nil {
		self.parser = data.CreateDataSource()
	}
	upload, _, err := request.FormFile("file")
	if err != nil {
		fmt.Fprint(w, "Error: a file is needed")
		return
	}
	out := csv.NewWriter(w)
	self.records = forecasting.PredictCSVSingle(upload)
	labels := make([]string, 2)
	labels[0] = self.records.Labels[0]
	labels[1] = self.records.Labels[5]
	err = out.Write(self.records.Labels)
	for i := 0; i<len(self.records.Data); i++ {
		out.Write(recordToString(self.records.Data[i]))
	}
	out.Flush()
	return
}
	
