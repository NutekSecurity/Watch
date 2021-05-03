package Wiki

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Deaths struct {
	Date      string `json:"date"`
	Wikipedia string `json:"wikipedia"`
	Deaths    []struct {
		Year        int    `json:"year"`
		Description string `json:"description"`
		Wikipedia   []struct {
			Title     string `json:"title"`
			Wikipedia string `json:"wikipedia"`
		} `json:"wikipedia"`
	} `json:"deaths"`
}

func GetDeaths(day int, month int) *Deaths {
	if day > 31 || month > 12 {
		return &Deaths{}
	}
	baseUrl := "https://byabbe.se/on-this-day/"
	getUrl := baseUrl + fmt.Sprintf("%d/%d/deaths.json", month, day)
	resp, err := http.Get(getUrl)
	if err != nil {
		return &Deaths{}
	}
	deaths := &Deaths{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(deaths)
	return deaths
}
