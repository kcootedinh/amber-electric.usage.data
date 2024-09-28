package amber

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Usage struct {
	Type              string    `json:"type"`
	Duration          int32     `json:"duration"`
	Date              string    `json:"date"`
	EndTime           time.Time `json:"endTime"`
	Quality           string    `json:"quality"`
	Kwh               float64   `json:"kwh"`
	NemTime           time.Time `json:"nemTime"`
	PerKwh            float64   `json:"perKwh"`
	ChannelType       string    `json:"channelType"`
	ChannelIdentifier string    `json:"channelIdentifier"`
	Cost              float64   `json:"cost"`
	Renewables        float64   `json:"renewables"`
	SpotPerKwh        float64   `json:"spotPerKwh"`
	StartTime         time.Time `json:"startTime"`
	SpikeStatus       string    `json:"spikeStatus"`
	TariffInformation struct {
		DemandWindow bool `json:"demandWindow"`
	} `json:"tariffInformation"`
	Descriptor string `json:"descriptor"`
}

const resolution = "30"

func (a amber) GetUsage(startDate, endDate time.Time) ([]Usage, error) {
	usageURL := a.constructURL(startDate, endDate)
	req, err := createRequest(usageURL, a.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("error constructing usage request: %w", err)
	}

	res, err := a.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making usage request: %w", err)
	}
	defer res.Body.Close()

	return decodeResponse(res)
}

func (a amber) constructURL(startDate, endDate time.Time) string {
	params := url.Values{}
	params.Add("startDate", startDate.Format("2006-01-02"))
	params.Add("endDate", endDate.Format("2006-01-02"))
	params.Add("resolution", resolution)
	return fmt.Sprintf("%s/sites/%s/usage?%s", a.BaseUrl, a.Site, params.Encode())
}

func createRequest(usageURL, apiKey string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, usageURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	return req, nil
}

func decodeResponse(res *http.Response) ([]Usage, error) {
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("usage request unsuccessful: %s", res.Status)
	}
	var usages []Usage
	if err := json.NewDecoder(res.Body).Decode(&usages); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}
	return usages, nil
}
