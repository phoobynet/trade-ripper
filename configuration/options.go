package configuration

type Options struct {
	DBHost         string `arg:"required,-h,--host" help:"The questdb post e.g. my.questdb.db"`
	DBInfluxPort   int    `arg:"-i,--influx" help:"The questdb influx port e.g. 9009 (default)"`
	DBPostgresPort int    `arg:"-p,--postgres" help:"The questdb postgres port e.g. 8812 (default)"`
	Class          string `arg:"required,-c,--class" help:"The class to subscribe to, either crypto or us_equity"`
	WebServerPort  int    `arg:"-w,--webserver" help:"The webserver port e.g. 3000 (default)"`
	Indexes        string `arg:"--indexes" help:"sp500,nasdaq100 - Currently only sp500 and nasdaq100 are supported.  Limits the data to the indexes specified"`
}
