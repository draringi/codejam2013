package data

import (
	"encoding/csv"
	"time"
	"fmt
	"strconv"
)

func csvParse(file io.Reader) (labels []string, data []Record) {
	reader := csv.NewReader (file)
	tmpdata, err := reader.ReadAll()
	if  err != nil {
		fmt.Println(err)
	}
	labels := tmpdata[0]
	data := make([]Record, len(tmpdata)-1)
	for i := 1; i<len(tmpdata); i++ {
		data[i-1].Time := time.Parse(dataSource.ISO, tmpdata[i][0])
		data[i-1].Radiation, err := strconv.ParseFloat(tmpdata[i][1], 32)
		if err != nil {
			data[i-1].empty = true
		}
		data[i-1].Humidity, err := strconv.ParseFloat(tmpdata[i][2], 32)
		if err != nil {
			data[i-1].empty = true
		}
		data[i-1].Temperature, err := strconv.ParseFloat(tmpdata[i][2], 32)
		if err != nil {
			data[i-1].empty = true
		}
		data[i-1].Wind, err := strconv.ParseFloat(tmpdata[i][2], 32)
		if err != nil {
			data[i-1].empty = true
		}
		data[i-1].Power, err := strconv.ParseFloat(tmpdata[i][2], 32)
		if err != nil {
			data[i-1].null = true
		}
	}
}

func fillRecords (emptyData []Record) (data []Record){
	gradRad, gradHumidity, gradTemp, gradWind := 0.0, 0.0, 0.0, 0.0
	for i := 0; i<len(emptyData); i++ {
		if emptyData[i].empty {
			emptyData[i].Radiation = emptyData[i-1].Radiation + gradRad
			emptyData[i].Humidity = emptyData[i-1].Humidity + gradHumidity
			emptyData[i].Temperature = emptyData[i-1].Temperature + gradTemp
			emptyData[i].Wind = emptyData[i-1].Wind + gradWind
			emptyData[i].empty = false
		} else {
			gradRad = (emptyData[i+4].Radiation - emptyData[i].Radiation)/4
			gradHumidity = (emptyData[i+4].Humidity - emptyData[i].Humidity)/4
			gradTemp = (emptyData[i+4].Temperature - emptyData[i].Temperature)/4
			gradWind = (emptyData[i+4].Wind - emptyData[i].Wind)/4
		}
	}
}
