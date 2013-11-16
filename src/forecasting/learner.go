package forecasting

import (
	"github.com/fxsjy/RF.go/RF"
	"io"
	"../data"
)

func buildData (data []data.Record) (inputs [][]interface{} , targets []float64){
	for i := 0; i<len(data); i++ {
		if data[i].null {
			break
		}
		row := make([]interface{},5)
		row[0]=data[i].Time
		row[1]=data[i].Radiation
		row[2]=data[i].Humidity
		row[3]=data[i].Temperature
		row[4]=data[i].Wind
		inputs = append(inputs,row)
		targets = append(targets,data[i].Power)
	}
	return
}

func learnCSV (file io.Reader, channel chan *data.CSVRequest) *RF.Forest {
	ret := make(chan *data.CSVData, 1)
	request := new(data.CSVRequest)
	request.Return = ret
	request.Request = file
	channel <- request
	var resp data.CSVData
	for {
		resp <-ret
		if resp != nil {
			break
		}
	}
	inputs, targets := buildData(resp.Data)
	forest := RF.BuildForest(inputs, targets, len(targets), len(inputs),1)
	return forest
}
	
