package forecasting

import (
	"github.com/fxsjy/RF.go/RF"
	"io"
	"draringi/codejam2013/src/data"
	"strconv"
)

func buildDataToLearn (data []data.Record) (inputs [][]interface{} , targets []string){
	for i := 0; i< len(data); i++ {
		if data[i].Null {
			break
		}
		row := make([]interface{}, 7)
		row[0]=data[i].Time
		row[1]=data[i].Time.Day()
		row[2]=data[i].Time.Hour()
		row[3]=data[i].Radiation
		row[4]=data[i].Humidity
		row[5]=data[i].Temperature
		row[6]=data[i].Wind
		inputs = append(inputs,row)
		targets = append(targets,strconv.FormatFloat(data[i].Power,'f', -1, 64))
	}
	return
}

func learnCSV (file io.Reader, channel chan *data.CSVRequest) *RF.Forest {
	ret := make(chan (*data.CSVData), 1)
	request := new(data.CSVRequest)
	request.Return = ret
	request.Request = file
	channel <- request
	resp := new(data.CSVData)
	for {
		resp = <-ret
		if resp != nil {
			break
		}
	}
	inputs, targets := buildDataToLearn(resp.Data)
	forest := RF.BuildForest(inputs, targets, len(targets), len(inputs),1)
	return forest
}
	

func learnData (data []data.Record) *RF.Forest {
	inputs, targets := buildDataToLearn(data)
	forest := RF.BuildForest(inputs, targets, len(targets), len(inputs),1)
	return forest
}

