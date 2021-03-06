package forecasting

import (
	"net/http"
	"io"
//	"os"
	"database/sql"
	"draringi/codejam2013/src/data"
	"strconv"
	"time"
	"encoding/xml"
	"math"
)

const quarter = (15*time.Minute)
const apikey = "B25ECB703CD25A1423DC2B1CF8E6F008"
const day = "day"

func buildDataToGuess (data []data.Record) (inputs [][]interface{}){
	for i := 0; i< len(data) ; i++ {
		if data[i].Null {
			row := make([]interface{},7)
			row[0]=data[i].Time
			row[1]=data[i].Time.Day()
			row[2]=data[i].Time.Hour()
			row[3]=data[i].Radiation
			row[4]=data[i].Humidity
			row[5]=data[i].Temperature
			row[6]=data[i].Wind
			inputs = append(inputs,row)
		}
	}
	return
}

type bad struct {
}

func (bad *bad) Error () string {
	return "something bad happened"
}


func PredictCSV (file io.Reader, channel chan *data.CSVRequest) *data.CSVData {
	forest := learnCSV(file, channel)
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
	inputs := buildDataToGuess(resp.Data)
	var outputs []string
	for i := 0; i<len(inputs); i++ {
		outputs = append (outputs, forest.Predicate(inputs[i]))
	}
	k:=0
	for i := 0; i<len(resp.Data); i++ {
		if resp.Data[i].Null {
			resp.Data[i].Power, _ = strconv.ParseFloat(outputs[k], 64)
			k++
			resp.Data[i].Null = false
		}
	}
	return resp
}

func PredictCSVSingle (file io.Reader) (*data.CSVData, error) {
	resp := new(data.CSVData)
	var err error
	resp.Labels, resp.Data, err = data.CSVParse(file)
	if err != nil{
		return nil, err
	}
	forest := learnData( resp.Data )
	inputs := buildDataToGuess( resp.Data )
	var outputs []string
	for i := 0; i<len(inputs); i++ {
		outputs = append (outputs, forest.Predicate(inputs[i]))
	}
	solution := new(data.CSVData)
	solution.Labels = resp.Labels
	solution.Data = make([]data.Record, len(outputs))
	k:=0
	for i := 0; i<len(resp.Data); i++ {
		if resp.Data[i].Null {
			solution.Data[k].Time = resp.Data[i].Time
			solution.Data[k].Power, _ = strconv.ParseFloat(outputs[k], 64)
			k++
			resp.Data[i].Null = false
		}
	}
	return solution, nil
}

func stdDev (correct []data.Record, guessed []data.Record) (float64, error) {
	if len(correct) != len(guessed) {
		return 0, new(bad)
	}
	var res float64 = 0.0
	for i:= 0; i < len(correct); i++ {
		res = res + math.Abs(correct[i].Power - guessed[i].Power)
	}
	res = res / float64(len(correct))
	return res, nil
}

func GenSTDev (file io.Reader) (result float64) {
	_, Data, err := data.CSVParse(file)
	if err != nil {
		return -1
	}
	forest := learnData( Data )
	for i := 0; i < len( Data ); i++ {
		Data[i].Null = true
	}
	inputs := buildDataToGuess( Data )
	guess := Data
	for i := 0; i < len( Data ); i++ {
		guess[i].Power, _ = strconv.ParseFloat(forest.Predicate(inputs[i]), 64)
	}
	result, err = stdDev(Data, guess)
	if err != nil {
		return -1
	} else {
		return result
	}
}

const SQLTIME = "2006-01-02 15:04:05+00"

func getPastData() []data.Record {
	const db_provider = "postgres"

	var db, err = sql.Open(db_provider, data.DB_connection)
	if err != nil {
		panic(err)
	}
	defer func () {_ = db.Close()} ()
	records := make([]data.Record, 0)
	var rows *sql.Rows
	rows, err = db.Query("SELECT * FROM Records;")
	for rows.Next() {
		var record data.Record
		var id int
		var tempTime string
		err = rows.Scan(&id ,&tempTime, &record.Radiation, &record.Humidity, &record.Temperature, &record.Wind, &record.Power)
		if err != nil {
			panic(err)
		}
		record.Time, err = time.Parse(SQLTIME, tempTime)
		if err != nil {
			panic(err)
		}
		records = append(records, record)
	}
	return records
}

func getFuture (id int, duration string) (resp *http.Response, err error) {
	client := new(http.Client)
	request, err:= http.NewRequest("GET", "https://api.pulseenergy.com/pulse/1/points/"+strconv.Itoa(id)+"/data.xml?interval="+duration+"&start="+strconv.FormatInt(time.Now().Unix(),10), nil)
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


func getFutureData() []data.Record{

	resp, err := getFuture(66094, day) // Radiation
	if err != nil {
		panic(err)
	}
	RadList :=  parseXmlFloat64(resp.Body)
	resp.Body.Close()
	
	resp, err = getFuture(66095, day) // Humidity
	if err != nil {
		panic(err)
	}
	HumidityList := parseXmlFloat64(resp.Body)
	resp.Body.Close()

	resp, err = getFuture(66077, day) // Temperature
	if err != nil {
		panic(err)
	}
	TempList := parseXmlFloat64(resp.Body)
	resp.Body.Close()

	resp, err = getFuture(66096, day) // Wind
	if err != nil {
		panic(err)
	}
	WindList := parseXmlFloat64(resp.Body)
	resp.Body.Close()

	records := make([]data.Record, len(RadList)*4)
	for i := 0; i < len(records); i++ {
		records[i].Empty = true
		records[i].Null = true
	}
	for i := 0; i < len(RadList); i++ {
		var err error
		records[i*4].Time, err = time.Parse(data.ISO,RadList[i].Date)
		if err != nil { //If it isn't ISO time, it might be time since epoch, or ISO_LONG
			records[i*4].Time, err = time.Parse(data.ISO_LONG,RadList[i].Date)
			if err != nil {
				var tmp int64
				tmp, err = strconv.ParseInt(RadList[i].Date, 10, 64)
				if err != nil { //If it isn't an Integer, and isn't ISO time, I have no idea what's going on.
					panic (err)
				}
				records[i*4].Time = time.Unix(tmp,0)
			}
		}
		records[i*4].Radiation = RadList[i].Value
		records[i*4].Humidity = HumidityList[i].Value
		records[i*4].Temperature = TempList[i].Value
		records[i*4].Wind = WindList[i].Value
		records[i*4].Empty = false
	}
	return fillRecords(records)
}

func fillRecords (emptyData []data.Record) (data []data.Record){
	gradRad, gradHumidity, gradTemp, gradWind := 0.0, 0.0, 0.0, 0.0
	for i := 0; i<len(emptyData); i++ {
		if emptyData[i].Empty && i > 0 {
			emptyData[i].Radiation = emptyData[i-1].Radiation + gradRad
			emptyData[i].Humidity = emptyData[i-1].Humidity + gradHumidity
			emptyData[i].Temperature = emptyData[i-1].Temperature + gradTemp
			emptyData[i].Wind = emptyData[i-1].Wind + gradWind
			emptyData[i].Time = emptyData[i-1].Time.Add(quarter)
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

func PredictPulse (Data chan ([]data.Record))  {
	notify := data.Monitor()
	for {
		if <-notify {
			forest := learnData(getPastData())
			pred := getFutureData()
			rawData := buildDataToGuess(pred)
			for i := 0; i < len(pred); i++ {
				forecast := forest.Predicate(rawData[i])
				pred[i].Power, _ = strconv.ParseFloat(forecast, 64)
			}
			Data <- pred
		} 
	}
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
