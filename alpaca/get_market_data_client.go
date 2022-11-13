package alpaca

import (
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var marketDataClient *resty.Client

func getMarketDataClient() *resty.Client {
	if marketDataClient == nil {
		marketDataClient = resty.New()
		marketDataClient.SetDebug(false)
		marketDataClient.SetBaseURL("https://data.alpaca.markets/v2")
		marketDataClient.SetHeaders(map[string]string{
			"APCA-API-KEY-ID":     os.Getenv("APCA_API_KEY_ID"),
			"APCA-API-SECRET-KEY": os.Getenv("APCA_API_SECRET_KEY"),
		})
	}

	return marketDataClient
}

func GetMarketData[T any](url string, queryParams map[string]string) (T, error) {
	var data T
	response, err := getMarketDataClient().R().SetQueryParams(queryParams).SetResult(&data).Get(url)

	if response.StatusCode() != http.StatusOK {
		logrus.Fatalf("Error getting data from %s: %s", url, response.Status())
	}

	return data, err
}
