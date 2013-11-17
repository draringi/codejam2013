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

func Monitor () (chan bool) {
	msg := make(chan bool, 5)
	go func () {
		db_init()
		getPastUnit(month) //Initialize the db with the past month's data
		for {
			getPastUnit(day)
			msg <- true //tell Predicate to update
			time.Sleep(quarter) //wait for another 15 mins
		}
	} ()
	return msg
}

func db_init() {
	var db, err = sql.Open(db_provider, db_connection)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS Records (ID SERIAL PRIMARY KEY UNIQUE,Time TIMESTAMP WITH TIME ZONE UNIQUE NOT NULL, Radiation DOUBLE precision, Humidity DOUBLE precision, Temperature DOUBLE precision, Wind DOUBLE precision, Power DOUBLE precision);")
    _, err = db.Exec("DROP FUNCTION IF EXISTS merge_Radiation ( timestamp with time zone, double precision) ;DROP FUNCTION IF EXISTS merge_Humidity ( timestamp with time zone, double precision) ;DROP FUNCTION IF EXISTS merge_Wind ( timestamp with time zone, double precision) ;DROP FUNCTION IF EXISTS merge_Temperature ( timestamp with time zone, double precision) ;") //clean out the functions, in case they are broken
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
	
func creativeUpdate(data []data.Record) {
	var db, err = sql.Open(db_provider, db_connection)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("CREATE FUNCTION merge_db(key timestamp with time zone, rad DOUBLE precision, humid DOUBLE precision, temp DOUBLE precision, w DOUBLE precision, pow DOUBLE precision) RETURNS VOID AS $$ BEGIN LOOP UPDATE Records SET Radiation = rad, Humidity=humid, Temperature=temp, Wind=w, Power=pow WHERE Time = key; IF found THEN RETURN; END IF; BEGIN INSERT INTO Records(Time, Radiation, Humidity, Temperature, Wind, Power) VALUES (key, rad, humid, temp, w, pow); RETURN; EXCEPTION WHEN unique_violation THEN END; END LOOP; END; $$ LANGUAGE plpgsql;")
    statement, staterr := db.Prepare("SELECT merge_db($1, $2, $3, $4, $5, $6);")
    if staterr != nil {
        panic(err)
    }
    defer statement.Close()
	for i := 0; i < len(data); i++ {
		_, err = statement.Exec(data[i].Time, data[i].Radiation, data[i].Humidity, data[i].Temperature, data[i].Wind, data[i].Power)
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
	
	resp, err = getPast(50578, unit) // Power
	if err != nil {
		panic(err)
	}
	PowerList := parseXmlFloat64(resp.Body)
	resp.Body.Close()

	recordList := buildRecord(RadList, HumidityList, TempList, WindList, PowerList)

	creativeUpdate(recordList)

}

import "fmt"

func buildRecord (RadList, HumidityList, TempList, WindList, PowerList []record) []data.Record {
	mult = len(PowerList)/len(RadList)
	list := make([]data.Record,len(PowerList))
	fmt.Println(strconv.Itoa(mult)
	for i := 0; i < len(PowerList); i++ {
		list[i].Empty = true
		list[i].Power = PowerList[i].Value
	}
	for i := 0; i < len(RadList); i++ {
		var err error
		list[i*4].Time, err = time.Parse(data.ISO,RadList[i].Date)
		if err != nil { //If it isn't ISO time, it might be time since epoch
			var tmp int64
			tmp, err = strconv.ParseInt(RadList[i].Date, 10, 64)
			if err != nil { //If it isn't an Integer, and isn't ISO time, I have no idea what's going on.
				panic (err)
			}
			records[i].Time = time.Unix(tmp,0)
		}
		list[i*mult].Radiation = RadList[i].Value
		list[i*mult].Humidity = HumidityList[i].Value
		list[i*mult].Temperature = TempList[i].Value
		list[i*mult].Wind = WindList[i].Value
		list[i*mult].Empty = false
	}
	return data.FillRecords(list)
}

