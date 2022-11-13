package alpaca

import (
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var tradeClient *resty.Client

func getTradeClient() *resty.Client {
	if tradeClient == nil {
		tradeClient = resty.New()
		tradeClient.SetDebug(false)
		tradeClient.SetBaseURL("https://paper-api.alpaca.markets/v2")
		tradeClient.SetHeaders(map[string]string{
			"APCA-API-KEY-ID":     os.Getenv("APCA_API_KEY_ID"),
			"APCA-API-SECRET-KEY": os.Getenv("APCA_API_SECRET_KEY"),
		})
	}

	return tradeClient
}

func GetTradeData[T any](url string, queryParams map[string]string) (T, error) {
	var data T
	response, err := getTradeClient().R().SetQueryParams(queryParams).SetResult(&data).Get(url)

	if response.StatusCode() != http.StatusOK {
		logrus.Fatalf("Error getting data from %s: %s", url, response.Status())
	}

	return data, err
}
