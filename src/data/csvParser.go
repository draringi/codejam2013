package data

import (
	"encoding/csv"
	"time"
	"fmt"
	"strconv"
	"io"
)

func CSVParse(file io.Reader) (labels []string, data []Record) {
	labels, data = csvParse(file)
	return 
} 

func csvParse(file io.Reader) (labels []string, data []Record) {
	reader := csv.NewReader (file)
	tmpdata, err := reader.ReadAll()
	if  err != nil {
		fmt.Println(err)
	}
	fmt.Println(len(tmpdata) - 1)
	labels = make([]string, 6)
	//labels = tmpdata[0]
	data  = make([]Record,  len(tmpdata)-1)
	for i := 1; i<len(tmpdata)-1; i++ {
		data[i-1].Time, _ = time.Parse(ISO, tmpdata[i][0])
		data[i-1].Radiation, err = strconv.ParseFloat(tmpdata[i][1], 64)
		if err != nil {
			data[i-1].empty = true
		}
		data[i-1].Humidity, err = strconv.ParseFloat(tmpdata[i][2], 64)
		if err != nil {
			data[i-1].empty = true
		}
		data[i-1].Temperature, err = strconv.ParseFloat(tmpdata[i][2], 64)
		if err != nil {
			data[i-1].empty = true
		}
		data[i-1].Wind, err = strconv.ParseFloat(tmpdata[i][2], 64)
		if err != nil {
			data[i-1].empty = true
		}
		data[i-1].Power, err = strconv.ParseFloat(tmpdata[i][2], 64)
		if err != nil {
			data[i-1].Null = true
		}
	}
	fmt.Println(len(data))
	data = fillRecords (data)
	return
}

func fillRecords (emptyData []Record) (data []Record){
	gradRad, gradHumidity, gradTemp, gradWind := 0.0, 0.0, 0.0, 0.0
	for i := 0; i<len(emptyData); i++ {
		if emptyData[i].empty && i > 0 {
			emptyData[i].Radiation = emptyData[i-1].Radiation + gradRad
			emptyData[i].Humidity = emptyData[i-1].Humidity + gradHumidity
			emptyData[i].Temperature = emptyData[i-1].Temperature + gradTemp
			emptyData[i].Wind = emptyData[i-1].Wind + gradWind
			emptyData[i].empty = false
		} else {
			if i + 4 < len (emptyData) {
				gradRad = (emptyData[i+4].Radiation - emptyData[i].Radiation)/4
				gradHumidity = (emptyData[i+4].Humidity - emptyData[i].Humidity)/4
				gradTemp = (emptyData[i+4].Temperature - emptyData[i].Temperature)/4
				gradWind = (emptyData[i+4].Wind - emptyData[i].Wind)/4
			} else {
				gradRad = 0
				gradHumidity = 0
				gradTemp = 0
				gradWind = 0
			}
		}
	}
	return
}
