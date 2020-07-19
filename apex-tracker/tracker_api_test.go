package tracker

import (
	"fmt"
	"os"
	"testing"
)

func TestTrackerApi(t *testing.T) {
	os.Setenv("TRACKER_TOKEN", "TOKEN")

	response, err := GetStats("LUV_nellrun", "psn")
	if err != nil {
		t.Error("something went wrong")
	}

	if response.RankScore.DisplayValue != "3,572" {
		t.Error(fmt.Sprintf("got %s, expected %s", response.RankScore.DisplayValue, "3,572"))
	}

	if response.RankScore.Value != 3572.000000 {
		t.Error(fmt.Sprintf("got %f, expected %f", response.RankScore.Value, 3572.000000))
	}
}
