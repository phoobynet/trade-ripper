package alpaca

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
)

var client *resty.Client

func getClient() *resty.Client {
	if client == nil {
		client = resty.New()
		client.SetDebug(true)
		client.SetBaseURL("https://paper-api.alpaca.markets/v2")
		client.SetHeaders(map[string]string{
			"APCA-API-KEY-ID":     os.Getenv("APCA_API_KEY_ID"),
			"APCA-API-SECRET-KEY": os.Getenv("APCA_API_SECRET_KEY"),
		})
	}

	return client
}

func GetData[T any](url string, queryParams map[string]string) (T, error) {
	fmt.Println("Getting data from: ", url)
	fmt.Printf("Query params: %v\n", queryParams)
	var data T
	_, err := getClient().R().SetQueryParams(queryParams).SetResult(&data).Get(url)

	return data, err
}
