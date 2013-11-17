package data

import (
	"encoding/csv"
	"time"
	"strconv"
	"io"
)

func CSVParse(file io.Reader) (labels []string, data []Record, err error) {
	labels, data, err = csvParseByLine(file)
	return 
} 

func csvParse(file io.Reader) (labels []string, data []Record, err error) {
	reader := csv.NewReader (file)
	tmpdata, err := reader.ReadAll()
	if  err != nil {
		return nil, nil, err
	}
	labels = tmpdata[0]
	data  = make([]Record,  len(tmpdata))
	for i := 1; i<len(tmpdata); i++ {
		data[i-1].Time, _ = time.Parse(ISO, tmpdata[i][0])
		data[i-1].Radiation, err = strconv.ParseFloat(tmpdata[i][1], 64)
		if err != nil {
			data[i-1].Empty = true
		}
		data[i-1].Humidity, err = strconv.ParseFloat(tmpdata[i][2], 64)
		if err != nil {
			data[i-1].Empty = true
		}
		data[i-1].Temperature, err = strconv.ParseFloat(tmpdata[i][3], 64)
		if err != nil {
			data[i-1].Empty = true
		}
		data[i-1].Wind, err = strconv.ParseFloat(tmpdata[i][4], 64)
		if err != nil {
			data[i-1].Empty = true
		}
		data[i-1].Power, err = strconv.ParseFloat(tmpdata[i][5], 64)
		if err != nil {
			data[i-1].Null = true
		}
	}
	data = FillRecords (data)
	return 
}

func csvParseByLine(file io.Reader) (labels []string, data []Record, err error) {
	reader := csv.NewReader (file)
	tmpdata, err := reader.Read()
	if  err != nil {
		return nil, nil, err
	}
	labels = tmpdata
	data  = make([]Record, 0)
	for {
		tmpdata, err := reader.Read()
		var rec Record
		if err == io.EOF {
			break
		} else if err == nil {
			rec.Time, err = time.Parse(ISO, tmpdata[0])
			if err != nil {
				break
			}
			rec.Radiation, err = strconv.ParseFloat(tmpdata[1], 64)
			if err != nil {
				rec.Empty = true
			}
			rec.Humidity, err = strconv.ParseFloat(tmpdata[2], 64)
			if err != nil {
				rec.Empty = true
			}
			rec.Temperature, err = strconv.ParseFloat(tmpdata[3], 64)
			if err != nil {
				rec.Empty = true
			}
			rec.Wind, err = strconv.ParseFloat(tmpdata[4], 64)
			if err != nil {
				rec.Empty = true
			}
			rec.Power, err = strconv.ParseFloat(tmpdata[5], 64)
			if err != nil {
				rec.Null = true
			}
		}
	}
	data = FillRecords (data)
	return 
}

func FillRecords (emptyData []Record) (data []Record){
	gradRad, gradHumidity, gradTemp, gradWind := 0.0, 0.0, 0.0, 0.0
	for i := 0; i<len(emptyData); i++ {
		if emptyData[i].Empty && i > 0 {
			emptyData[i].Radiation = emptyData[i-1].Radiation + gradRad
			emptyData[i].Humidity = emptyData[i-1].Humidity + gradHumidity
			emptyData[i].Temperature = emptyData[i-1].Temperature + gradTemp
			emptyData[i].Wind = emptyData[i-1].Wind + gradWind
			emptyData[i].Empty = false
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
	return emptyData
}
