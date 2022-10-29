package alpaca

import "net/url"

var usEquitiesURL url.URL = url.URL{Scheme: "wss", Host: "stream.data.alpaca.markets", Path: "/v2/sip"}
var cryptoURL url.URL = url.URL{Scheme: "wss", Host: "stream.data.alpaca.markets", Path: "/v1beta2/crypto"}
