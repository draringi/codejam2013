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
var Conf float64 = 1000

type MachineInterface struct {
	records *data.CSVData
	parser *data.DataSource
}

type Feeder struct {
	parser *data.DataSource
}

type StdDev struct {
	parser *data.DataSource
}

func (self *StdDev) ServeHTTP (w http.ResponseWriter, request *http.Request) {
	upload, _, err := request.FormFile("file")
	if err != nil {
		fmt.Fprint(w, "Error: a file is needed")
		return
	}
	stddev := forecasting.GenSTDev(upload)
	fmt.Fprint(w, strconv.FormatFloat(stddev,'f', -1, 64))
	return
}

func (self *Feeder) ServeHTTP (w http.ResponseWriter, request *http.Request) {
	if self.parser == nil {
		self.parser = data.CreateDataSource()
	}
	upload, _, err := request.FormFile("file")
	if err != nil {
		fmt.Fprint(w, "Error: a file is needed")
		return
	}
	data.AddCSVToDB(upload)
	fmt.Fprint(w, "Done")
	return
}

func recordToString(record data.Record) []string{
	stringRecord := make([]string, 2)
	stringRecord[0] = record.Time.Format(data.ISO)
	stringRecord[1] = strconv.FormatFloat(record.Power,'f', -1, 64)
	return stringRecord
}

func confToString(record data.Record) []string{
	stringRecord := make([]string, 3)
	stringRecord[0] = record.Time.Format(data.ISO)
	stringRecord[1] = strconv.FormatFloat(record.Power - Conf,'f', -1, 64)
	stringRecord[2] = strconv.FormatFloat(record.Power + Conf,'f', -1, 64)
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
	self.records, err = forecasting.PredictCSVSingle(upload)
	if err != nil {
		fmt.Fprint(w, "An Error Occured")
		return
	}
	labels := make([]string, 2)
	labels[0] = self.records.Labels[0]
	labels[1] = self.records.Labels[5]
	err = out.Write(labels)
	if err != nil {
		fmt.Fprint(w, "An Error Occured")
		return
	}
	for i := 0; i<len(self.records.Data); i++ {
		out.Write(recordToString(self.records.Data[i]))
	}
	out.Flush()
	return
}

type Confidence struct {
	records *data.CSVData
	parser *data.DataSource
}

func (self *Confidence) ServeHTTP (w http.ResponseWriter, request *http.Request) {
	if self.parser == nil {
		self.parser = data.CreateDataSource()
	}
	upload, _, err := request.FormFile("file")
	if err != nil {
		fmt.Fprint(w, "Error: a file is needed")
		return
	}
	out := csv.NewWriter(w)
	self.records, err = forecasting.PredictCSVSingle(upload)
	if err != nil {
		fmt.Fprint(w, "An Error Occured")
		return
	}
	labels := make([]string, 3)
	labels[0] = "Date"
	labels[1] = "Lower 75% Conf"
	labels[2] = "Upper 75% Conf"
	err = out.Write(labels)
	for i := 0; i<len(self.records.Data); i++ {
		out.Write(confToString(self.records.Data[i]))
	}
	out.Flush()
	return
}

