package alpaca

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"time"
)

// SIPReader - https://alpaca.markets/docs/api-references/market-data-api/stock-pricing-data/realtime/
type SIPReader struct {
	conn                 *websocket.Conn
	closed               bool
	config               *SIPReaderConfig
	restartCountInPeriod int
	lastRestartTime      time.Time
	lastMessageAt        time.Time
}

func NewSIPReader(config *SIPReaderConfig) *SIPReader {
	if config.SocketURL == nil {
		config.SocketURL = &sipSocketURL
	}
	return &SIPReader{
		config: config,
	}
}

func (r *SIPReader) Stop() error {
	if r.conn != nil && !r.closed {
		return r.conn.Close()
	}

	return nil
}

var restarts int

func (r *SIPReader) Start() error {
	defer func() {
		if recErr := recover(); recErr != nil {
			if restarts > 50 {
				panic(fmt.Errorf("too many restarts: %v", recErr))
			}

			logrus.Errorf("recovering from panic (will restart in 2 seconds): %v", recErr)
			time.Sleep(2 * time.Second)
			startErr := r.Start()

			if startErr != nil {
				panic(startErr)
			}

			restarts++
		}
	}()

	stopErr := r.Stop()

	if stopErr != nil {
		return stopErr
	}

	r.lastRestartTime = time.Now()

	logrus.Infof("connecting to socket @%s", r.config.SocketURL.String())

	conn, _, dialErr := websocket.DefaultDialer.Dial(r.config.SocketURL.String(), nil)

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

func (r *SIPReader) auth() error {
	logrus.Info("authenticating with alpaca...")
	return r.conn.WriteJSON(&authMessage{
		Action: "auth",
		Key:    r.config.Key,
		Secret: r.config.Secret,
	})
}

func (r *SIPReader) subscribe() error {
	logrus.Info("subscribing to trades...")
	return r.conn.WriteJSON(&subscribeMessage{
		Action: "subscribe",
		Trades: r.config.Trades,
	})
}
