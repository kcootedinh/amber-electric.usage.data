package amber

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type amber struct {
	BaseUrl string
	ApiKey  string
	Site    string
	Client  http.Client
	Limiter *rate.Limiter
}

type Service interface {
	GetUsage(ctx context.Context, startDate, endDate time.Time) ([]Usage, error)
}

func NewUsageService(baseUrl, apiKey, site string) Service {
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	limiter := rate.NewLimiter(rate.Every(300*time.Second), 50) // 50 request every 300 seconds
	return amber{BaseUrl: baseUrl, ApiKey: apiKey, Site: site, Client: client, Limiter: limiter}
}
