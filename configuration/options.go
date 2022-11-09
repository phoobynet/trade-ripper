package configuration

type Options struct {
	DBHost         string `arg:"required,-h,--host" help:"The questdb post e.g. my.questdb.db"`
	DBInfluxPort   int    `arg:"-i,--influx" help:"The questdb influx port e.g. 9009 (default)"`
	DBPostgresPort int    `arg:"-p,--postgres" help:"The questdb postgres port e.g. 8812 (default)"`
	Class          string `arg:"required,-c,--class" help:"The class to subscribe to, either crypto or us_equity"`
	WebServerPort  int    `arg:"-w,--webserver" help:"The webserver port e.g. 3000 (default)"`
	Indexes        string `arg:"--indexes" help:"example: sp500,nasdaq100,djia - Currently only sp500, nasdaq100, djia are supported.  Limits the data to the indexes specified"`
	TickersFile    string `arg:"--tickersfile" help:"example: tickers.txt - A file containing a list of tickers to subscribe to, seperated by newlines.  Can be used in conjunction with --indexes"`
	Tickers        string `arg:"--tickers" help:"example: AAPL,MSFT,GOOG - A comma seperated list of tickers to subscribe to.  Can be used in conjunction with --indexes and --tickersfile"`
}
