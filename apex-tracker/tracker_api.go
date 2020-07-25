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
	SegmentType string   `json:"type"`
	Metadata    Metadata `json:"metadata"`
	Stats       Stats    `json:"stats"`
}

// Metadata meta
type Metadata struct {
	Name         string `json:"name"`
	ImageURL     string `json:"imageUrl"`
	TallImageURL string `json:"tallImageUrl"`
}

// Stats one of player stats
type Stats struct {
	RankScore Stat `json:"rankScore"`
	Kills     Stat `json:"kills"`
	Damage    Stat `json:"damage"`
}

// Stat describe ranked league score
type Stat struct {
	Value        float64      `json:"value"`
	DisplayValue string       `json:"displayValue"`
	Metadata     StatMetadata `json:"metadata"`
}

// StatMetadata metadata
type StatMetadata struct {
	RankName string `json:"rankName"`
}

// GetStats getting user stats from tracker
func GetStats(username string, platform string) ([]Segment, error) {
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

	return playserProfileResponse.Data.Segments, nil
}
