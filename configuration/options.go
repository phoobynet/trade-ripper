package configuration

type Options struct {
	QuestDBURI    string `arg:"required,-q" help:"The questdb post e.g. my.questdb.db:9009"`
	WebServerPort int    `arg:"-p" help:"The port to run the web server on"`
}
