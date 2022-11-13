package alpaca

import (
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var client *resty.Client

func getClient() *resty.Client {
	if client == nil {
		client = resty.New()
		client.SetDebug(false)
		client.SetBaseURL("https://paper-api.alpaca.markets/v2")
		client.SetHeaders(map[string]string{
			"APCA-API-KEY-ID":     os.Getenv("APCA_API_KEY_ID"),
			"APCA-API-SECRET-KEY": os.Getenv("APCA_API_SECRET_KEY"),
		})
	}

	return client
}

func GetData[T any](url string, queryParams map[string]string) (T, error) {
	var data T
	response, err := getClient().R().SetQueryParams(queryParams).SetResult(&data).Get(url)

	if response.StatusCode() != http.StatusOK {
		logrus.Fatalf("Error getting data from %s: %s", url, response.Status())
	}

	return data, err
}
