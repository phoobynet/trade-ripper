package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"github.com/phoobynet/trade-ripper/configuration"
	"github.com/sirupsen/logrus"
)

type client struct {
	conn *websocket.Conn
}

var connectedClients = make(map[string]*client)

type Message struct {
	MessageType    string `json:"type"`
	MessageContent string `json:"content"`
}

func Broadcast(message Message) {
	for _, c := range connectedClients {
		err := c.conn.WriteJSON(message)

		if err != nil {
			logrus.Errorf("Error writing message to client: %v", err)
		}
	}
}

func Run(options configuration.Options) {
	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	app.Static("/", "./public")

	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			clientID := uuid.NewString()
			connectedClients[clientID] = nil
			c.Locals("clientID", clientID)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		clientID := c.Locals("clientID").(string)

		if _, ok := connectedClients[clientID]; ok {
			logrus.Infof("Client %s connected", clientID)
			c.SetCloseHandler(func(code int, text string) error {
				logrus.Infof("Client %s disconnected", clientID)
				delete(connectedClients, clientID)
				return nil
			})

			connectedClients[clientID] = &client{conn: c}
		}

		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				logrus.Panicf("Error reading message: %v", err)
				break
			}

			// TODO: What do I do with the message?
			logrus.Infof("recv: %s (%d)", msg, mt)
		}
	}))

	logrus.Fatalln(app.Listen(fmt.Sprintf(":%d", options.WebServerPort)))
}
