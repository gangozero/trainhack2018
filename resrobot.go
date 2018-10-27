package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Board struct {
	Arrival []BoardItem `json:"Arrival"`
}

type BoardItem struct {
	Time            string `json:"time"`
	Date            string `json:"date"`
	TransportNumber string `json:"transportNumber"`
}

// get time or selected train from ResRobot
func getTime(stationID string, trainID string) (time.Time, error) {
	url := "https://api.resrobot.se/v2/arrivalBoard"
	key := os.Getenv("TRAFIKLAB_TIMETABLE_KEY")

	params := map[string]string{
		"key":         key,
		"id":          stationID,
		"maxJourneys": "50",
		"format":      "json",
		"operators":   "74",
		"passlist":    "0",
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return time.Now(), fmt.Errorf("Error creating HTTP request: %s", err.Error())
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}

	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return time.Now(), fmt.Errorf("Error getting HTTP data: %s", err.Error())
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return time.Now(), fmt.Errorf("Error reading response body: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return time.Now(), fmt.Errorf("Wrong response code '%d' with body '%s'", resp.StatusCode, string(bodyBytes))
	}

	var board Board
	err = json.Unmarshal(bodyBytes, &board)
	if err != nil {
		return time.Now(), fmt.Errorf("Error parsing response json: %s", err.Error())
	}

	for _, b := range board.Arrival {
		if b.TransportNumber == trainID {
			joinTime := fmt.Sprintf("%sT%s-02:00", b.Date, b.Time)
			tp, err := time.Parse(time.RFC3339, joinTime)
			if err != nil {
				return time.Now(), fmt.Errorf("Error parsing time: %s", err.Error())
			}
			return tp, nil
		}
	}

	return time.Now(), fmt.Errorf("Train %s not found in response: %v", trainID, board)
}
