package alpaca

import "testing"

func BenchmarkConvertToTradeRows(b *testing.B) {
	inputJSON := []byte(`[{"T": "t", "S": "AAPL", "s": 100, "p": 100, "t": "2021-01-01T00:00:00.00000Z", "c": ["@"], "z":"A", "x": "Z" }]`)

	for i := 0; i < b.N; i++ {
		_, err := ConvertToTradeRows(inputJSON)
		if err != nil {
			b.Errorf("failed to convert to trade rows: %v", err)
		}
	}
}
