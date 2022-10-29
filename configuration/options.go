package configuration

type Options struct {
	QuestDBURI string `arg:"required,-q" help:"The questdb post e.g. my.questdb.db:9009"`
	Class      string `arg:"required,-c" help:"The class to subscribe to, either crypto or us_equity"`
}
