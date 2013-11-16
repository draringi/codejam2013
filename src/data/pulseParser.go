package data

import (
	"net/http"
	"encoding/xml"
	"io"
	"strconv"
	"time"
	"database/sql"
	"os"
	_ "github.com/jbarham/gopgsqldriver"
)

var db_connection = "user=adminficeuc6 dbname=codejam2013 password=zUSfsRCcvNZf host="+os.Getenv("OPENSHIFT_POSTGRESQL_DB_HOST")+" port="+os.Getenv("OPENSHIFT_POSTGRESQL_DB_PORT")
const db_provider = "postgres"

const apikey = "B25ECB703CD25A1423DC2B1CF8E6F008"

const day = "day"
const month = "month"

const quarter = (15*time.Minute)

func Monitor (msg chan bool) {
	
	go func () {
		msg = make(chan bool, 1)
		db_init()
		getPastUnit(month) //Initialize the db with the past month's data
		for {
			getPastUnit(day)
			msg <- true //tell Predicate to update
			time.Sleep(quarter) //wait for another 15 mins
		}
	} ()
	return
}

func db_init() {
	var db, err = sql.Open(db_provider, db_connection)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS Records (ID SERIAL PRIMARY KEY UNIQUE,Time TIMESTAMP WITH TIME ZONE UNIQUE NOT NULL, Radiation DOUBLE precision, Humidity DOUBLE precision, Temperature DOUBLE precision, Wind DOUBLE precision, Power DOUBLE precision);")
}

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

func parseXmlFloat64 (r io.Reader) []record {
	decoder := xml.NewDecoder(r)
	var output point
	err := decoder.Decode(&output)
	if err != nil {
		panic(err)
	}
	return output.Records.RecordList
}
	
func creativeUpdate(field string, data []record) {
	var db, err = sql.Open(db_provider, db_connection)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("CREATE FUNCTION merge_db(key timestamp with time zone, data DOUBLE precision) RETURNS VOID AS $$ BEGIN LOOP UPDATE db SET $1 = data WHERE Time = key; IF found THEN RETURN; END IF; BEGIN INSERT INTO Records(Time,$1) VALUES (key, data); RETURN; EXCEPTION WHEN unique_violation THEN END; END LOOP; END; $$ LANGUAGE plpgsql;", field)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(data); i++ {
		_, err = db.Exec("SELECT merge_db($1, $2);", data[i].Date, data[i].Value)
		if err != nil {
			panic(err)
		}
	}
}


func getPastUnit (unit string) {
	resp, err := getPast(66094, unit) // Radiation
	if err != nil {
		panic(err)
	}
	RadList :=  parseXmlFloat64(resp.Body)
	resp.Body.Close()
	
	resp, err = getPast(66095, unit) // Humidity
	if err != nil {
		panic(err)
	}
	HumidityList := parseXmlFloat64(resp.Body)
	resp.Body.Close()

	resp, err = getPast(66077, unit) // Temperature
	if err != nil {
		panic(err)
	}
	TempList := parseXmlFloat64(resp.Body)
	resp.Body.Close()

	resp, err = getPast(66096, unit) // Wind
	if err != nil {
		panic(err)
	}
	WindList := parseXmlFloat64(resp.Body)
	resp.Body.Close()
	
	resp, err = getPast(66095, unit) // Power
	if err != nil {
		panic(err)
	}
	PowerList := parseXmlFloat64(resp.Body)
	resp.Body.Close()

	creativeUpdate("Radiation", RadList)
	creativeUpdate("Humidity", HumidityList)
	creativeUpdate("Temperature", TempList)
	creativeUpdate("Wind", WindList)
	creativeUpdate("Power", PowerList)
}

