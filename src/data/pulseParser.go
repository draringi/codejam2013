package data

import (
	"net/http"
	"encoding/json"
	"io"
	"strconv"
)

const apikey = "B25ECB703CD25A1423DC2B1CF8E6F008"

const day = "day"

func getPast (id int, duration string) (resp *http.Response, err error) {
	client := new(http.Client)
	request, err:= http.NewRequest("GET", "https://api.pulseenergy.com/pulse/1/points/"+strconv.Itoa(id)+"/data.json?interval="+duration, nil)
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

/*
typedef floattimelist struct {
	time string
	
}

typedef jsonfloat struct {
	id int
	label, unit, quantity, resource, start, end string
	average float64
}

func parseJsonFloat64 (r io.Reader) (
	decoder = json.NewDecoder(r)
	

func getPastDay () ([]Record) {
	resp, err := getPast(66094, day) // Radiation
	if err != nil {
		//Something bad
	}
	RadList := 
	
*/
