package amber

import (
	"net/http"
	"time"
)

type amber struct {
	BaseUrl string
	ApiKey  string
	Site    string
	Client  http.Client
}

type Service interface {
	GetUsage(startDate, endDate time.Time) ([]Usage, error)
}

func NewUsageService(baseUrl, apiKey, site string) Service {
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	return amber{BaseUrl: baseUrl, ApiKey: apiKey, Site: site, Client: client}
}
