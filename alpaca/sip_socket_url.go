package alpaca

import "net/url"

var sipSocketURL = url.URL{Scheme: "wss", Host: "stream.data.alpaca.markets", Path: "/v2/sip"}
