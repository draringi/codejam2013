package forecasting

import (
	"io"
	"../data"
	"strconv"
)

func buildDataToGuess (data []data.Record) (inputs [][]interface{}){
	for i := 0; i<len(data); i++ {
		if data[i].Null {
			row := make([]interface{},5)
			row[0]=data[i].Time
			row[1]=data[i].Radiation
			row[2]=data[i].Humidity
			row[3]=data[i].Temperature
			row[4]=data[i].Wind
			inputs = append(inputs,row)
		}
	}
	return
}

func PredictCSV (file io.Reader, channel chan *data.CSVRequest) []data.Record {
	forest := learnCSV(file, channel)
	ret := make(chan (*data.CSVData), 1)
	request := new(data.CSVRequest)
	request.Return = ret
	request.Request = file
	channel <- request
	var resp *data.CSVData
	for {
		resp = <-ret
		if resp != nil {
			break
		}
	}
	inputs := buildDataToGuess(resp.Data)
	var outputs []string
	for i := 0; i<len(inputs); i++ {
		outputs = append (outputs, forest.Predicate(inputs[i]))
	}
	k:=0
	for i = 0; i<len(resp.Data); i++ {
		if resp.Data[i].Null {
			resp.Data[i].Power, _ = strconv.ParseFloat(outputs[k], 64)
			k++
			resp.Data[i].Null = false
		}
	}
	return resp.Data
}


