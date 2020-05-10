package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type ta struct {
	TimeFrame int    `json:"timeFrame"`
	Type      string `json:"type"`
}

// Request represents an alpha request
type Request struct {
	OrderbookID          int       `json:"orderbookId"`
	ChartType            string    `json:"chartType"`
	WidthOfPlotContainer int       `json:"widthOfPlotContainer"`
	ChartResolution      string    `json:"chartResolution"`
	Navigator            bool      `json:"navigator"`
	Percentage           bool      `json:"percentage"`
	Volume               bool      `json:"volume"`
	Owners               bool      `json:"owners"`
	Start                time.Time `json:"start"`
	End                  time.Time `json:"end"`
	Ta                   []ta      `json:"ta"`
	CompareIds           []int     `json:"compareIds"`
}

// Response represents an alpha response
type Response struct {
	DataPoints         [][]float64 `json:"dataPoints"`
	TrendSeries        [][]float64 `json:"trendSeries"`
	AllowedResolutions []string    `json:"allowedResolutions"`
	DefaultResolution  string      `json:"defaultResolution"`
	Comparisons        []struct {
		OrderbookName string      `json:"orderbookName"`
		ShortName     string      `json:"shortName"`
		OrderbookID   string      `json:"orderbookId"`
		DataPoints    [][]float64 `json:"dataPoints"`
	} `json:"comparisons"`
	TechnicalAnalysis []struct {
		DataPoints [][]float64 `json:"dataPoints"`
		TimeFrame  int         `json:"timeFrame"`
		Type       string      `json:"type"`
	} `json:"technicalAnalysis"`
	OwnersPoints  []interface{} `json:"ownersPoints"`
	ChangePercent float64       `json:"changePercent"`
	High          float64       `json:"high"`
	LastPrice     float64       `json:"lastPrice"`
	Low           float64       `json:"low"`
}

type csvFormat struct {
	Time  float64
	Price float64
	Ma50  float64
	Ma200 float64
	Ema21 float64
}

func getRequest(orderbookID int) Request {
	var req Request
	req.OrderbookID = orderbookID
	req.ChartType = "AREA"
	req.WidthOfPlotContainer = 558
	req.ChartResolution = "DAY"
	req.Navigator = false
	req.Percentage = false
	req.Volume = false
	req.Owners = false
	req.Start, _ = time.Parse(time.RFC3339, "2015-01-01T22:00:00.000Z")
	req.End, _ = time.Parse(time.RFC3339, "2020-10-01T22:00:00.000Z")
	req.Ta = []ta{ta{50, "sma"}, ta{200, "sma"}, ta{21, "ema"}}
	req.CompareIds = []int{19002}

	return req
}

func getResponse(orderbookID int) Response {
	client := &http.Client{}

	reqBody, err := json.Marshal(getRequest(orderbookID))

	println(string(reqBody))
	req, err := http.NewRequest("POST", "https://limitless-solar-winds.herokuapp.com/https://www.avanza.se/ab/component/highstockchart/getchart/orderbook", bytes.NewBuffer(reqBody))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Access-Control-Allow-Origin", "*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Referer", "https://www.example.com")
	req.Header.Add("Origin", "https://www.example.com")

	resp, err := client.Do(req)

	if err != nil {
		println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println(err)
	}

	var response Response

	err = json.Unmarshal(body, &response)
	if err != nil {
		println(err)
	}

	return response
}

func main() {
	response := getResponse(599956)

	filename := fmt.Sprintf("%d.csv", 599956)

	err := os.Remove(filename)
	if os.IsNotExist(err) {
	}
	file, err := os.Create(filename)

	if err != nil {
		println(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	for i := 0; i < len(response.DataPoints); i++ {
		if len(response.TechnicalAnalysis[0].DataPoints) <= i {
			continue
		} else if len(response.TechnicalAnalysis[1].DataPoints) <= i {
			continue
		} else if len(response.TechnicalAnalysis[2].DataPoints) <= i {
			continue
		}

		timestamp := fmt.Sprintf("%f", response.DataPoints[i][0])
		price := fmt.Sprintf("%f", response.DataPoints[i][1])
		ma50 := fmt.Sprintf("%f", response.TechnicalAnalysis[0].DataPoints[i][1])
		ma200 := fmt.Sprintf("%f", response.TechnicalAnalysis[1].DataPoints[i][1])
		ema21 := fmt.Sprintf("%f", response.TechnicalAnalysis[2].DataPoints[i][1])

		err := writer.Write([]string{timestamp, price, ma50, ma200, ema21})
		if err != nil {
			println(err)
		}
	}
}
