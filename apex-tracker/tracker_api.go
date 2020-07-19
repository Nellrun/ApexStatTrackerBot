package tracker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// PlayserProfileResponse main body
type PlayserProfileResponse struct {
	Data Data `json:"data"`
}

// Data Main response body
type Data struct {
	Segments []Segment `json:"segments"`
}

// Segment describe one segment
type Segment struct {
	SegmentType string `json:"type"`
	Stats       Stats  `json:"stats"`
}

// Stats one of player stats
type Stats struct {
	RankScore RankScore `json:"rankScore"`
}

// RankScore describe ranked league score
type RankScore struct {
	Value        float64 `json:"value"`
	DisplayValue string  `json:"displayValue"`
}

// GetStats getting user stats from tracker
func GetStats(username string, platform string) (*Stats, error) {
	client := &http.Client{}

	url := fmt.Sprintf("https://public-api.tracker.gg/v2/apex/standard/profile/%s/%s", platform, username)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	request.Header.Add("TRN-Api-Key", os.Getenv("TRACKER_TOKEN"))

	response, err := client.Do(request)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	var playserProfileResponse PlayserProfileResponse

	err = json.Unmarshal(body, &playserProfileResponse)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return &playserProfileResponse.Data.Segments[0].Stats, nil
}
