package adapters

import (
	"testing"
)

func Test_AdaptRawMessageToTrades(t *testing.T) {
	rawMessage := []byte(`[{
  "T": "t",
  "i": 96921,
  "S": "AAPL",
  "x": "D",
  "p": 126.55,
  "s": 1,
  "t": "2021-02-22T15:51:44.208Z",
  "c": ["@", "I"],
  "z": "C"
}]`)

	actual, err := AdaptRawMessageToTrades(rawMessage)

	if err != nil {
		t.Errorf("failed to adapter message: %v", err)
	}

	if len(actual) != 1 {
		t.Errorf("expected 1 trade, got %v", len(actual))
	}

	t.Log(actual)

	if actual[0]["S"] != "AAPL" {
		t.Errorf("expected symbol to be AAPL, got %v", actual[0]["S"])
	}
}
