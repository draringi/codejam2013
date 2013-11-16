package data

import (
	"net/http"
	"encoding/xml"
//	"io"
	"strconv"
	"time"
//	"database/sql"
)

//import _ "github.com/jbarham/gopgsqldriver"

const apikey = "B25ECB703CD25A1423DC2B1CF8E6F008"

const day = "day"

const quarter = (15*time.Minute)

func getPast (id int, duration string) (resp *http.Response, err error) {
	client := new(http.Client)
	request, err:= http.NewRequest("GET", "https://api.pulseenergy.com/pulse/1/points/"+strconv.Itoa(id)+"/data.xml?interval="+duration, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", apikey)
	resp, err = client.Do(request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type records struct {
	RecordList []record `xml:"record"`
}

type record struct {
	Date string `xml:"date,attr"`
	Value float64 `xml:"value,attr"`
}

type point struct {
	Records records `xml:"records"`
}

func parseXmlFloat64 (r io.Reader) [][]record(
	decoder = xml.NewDecoder(r)
	var output point
	err := decoder.decode(&output)
	return output.Records.RecordList
	

func getPastDay () {
	resp, err := getPast(66094, day) // Radiation
	if err != nil {
		panic(err)
	}
	RadList :=  parseXmlFloat64(resp.Body)
	resp.Body.Close()
	
	resp, err = getPast(66095, day) // Humidity
	if err != nil {
		panic(err)
	}
	HumidityList := parseXmlFloat64(resp.Body)
	resp.Body.Close()

	resp, err = getPast(66077, day) // Temp
	if err != nil {
		panic(err)
	}
	HumidityList := parseXmlFloat64(resp.Body)
	resp.Body.Close()

	resp, err = getPast(66096, day) // Wind
	if err != nil {
		panic(err)
	}
	HumidityList := parseXmlFloat64(resp.Body)
	resp.Body.Close()
	
	resp, err = getPast(66095, day) // Power
	if err != nil {
		panic(err)
	}
	HumidityList := parseXmlFloat64(resp.Body)
	resp.Body.Close()

