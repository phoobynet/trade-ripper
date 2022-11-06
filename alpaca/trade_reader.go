package alpaca

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/url"
	"time"
)

// TradeReader - https://alpaca.markets/docs/api-references/market-data-api/stock-pricing-data/realtime/
type TradeReader struct {
	conn            *websocket.Conn
	closed          bool
	config          *TradeReaderConfig
	lastRestartTime time.Time
	lastMessageAt   time.Time
	socketURL       url.URL
}

func NewTradeReader(config *TradeReaderConfig) *TradeReader {
	var socketURL url.URL

	if config.Options.Class == "crypto" {
		socketURL = cryptoURL
	} else if config.Options.Class == "us_equity" {
		socketURL = usEquitiesURL
	}

	return &TradeReader{
		config:    config,
		socketURL: socketURL,
	}
}

func (r *TradeReader) Stop() error {
	if r.conn != nil && !r.closed {
		return r.conn.Close()
	}

	return nil
}

func (r *TradeReader) Start() error {
	defer func() {
		if rec := recover(); rec != nil {
			r.config.webServer.Publish(map[string]any{
				"type":    "error",
				"message": fmt.Sprintf("panic and recovered: %v", rec),
			})

			logrus.Error("recovering from panic: ", rec)
			time.Sleep(2 * time.Second)
			restartErr := r.Start()

			if restartErr != nil {
				panic(fmt.Errorf("failed to restart: %w", restartErr))
			}
		}
	}()
	stopErr := r.Stop()

	if stopErr != nil {
		return stopErr
	}

	r.lastRestartTime = time.Now()

	logrus.Infof("connecting to socket @%s", r.socketURL.String())

	conn, _, dialErr := websocket.DefaultDialer.Dial(r.socketURL.String(), nil)

	if dialErr != nil {
		return dialErr
	}

	r.closed = false
	r.conn = conn

	r.conn.SetCloseHandler(func(code int, text string) error {
		logrus.Warnf("connection closed: %d %s", code, text)
		r.closed = true
		return nil
	})

	authErr := r.auth()

	if authErr != nil {
		return authErr
	}

	subscribeErr := r.subscribe()

	if subscribeErr != nil {
		return subscribeErr
	}

	for {
		_, rawMessage, readMessageError := r.conn.ReadMessage()

		if readMessageError != nil {
			r.config.ErrorsChannel <- fmt.Errorf("read message failed: %+v", readMessageError)
			_ = conn.Close()
			r.closed = true
			panic(readMessageError)
		}

		r.config.RawMessageChannel <- rawMessage
		r.lastMessageAt = time.Now()
	}
}

func (r *TradeReader) auth() error {
	logrus.Info("authenticating with alpaca...")
	return r.conn.WriteJSON(&authMessage{
		Action: "auth",
		Key:    r.config.Key,
		Secret: r.config.Secret,
	})
}

func (r *TradeReader) subscribe() error {
	logrus.Info("subscribing to trades...")
	return r.conn.WriteJSON(&subscribeMessage{
		Action: "subscribe",
		Trades: r.config.Symbols,
	})
}
