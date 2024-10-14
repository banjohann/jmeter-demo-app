package main

import (
	server "github.com/JohannBandelow/jmeter-demo-app/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	viewsEngine := html.New("./web/static/views", ".html")

	app := fiber.New(fiber.Config{
		Views: viewsEngine,
	})

	appHandler := server.NewAppHandler()
	wsServer := server.NewWebSocket()

	app.Static("/static/", "./web/static")

	app.Get("/", appHandler.HandleGetIndex)

	app.Get("/ws", websocket.New(func(ctx *websocket.Conn) {
		wsServer.HandleWebSocket(ctx)
	}))

	app.Post("/message", func(c *fiber.Ctx) error {
		msg := new(server.Message)

		if err := c.BodyParser(msg); err != nil {
			return err
		}

		wsServer.PublishMessageAllClients(msg)

		return c.Send([]byte("OK"))
	})

	go wsServer.HandleMessages()

	app.Listen(":3000")
}
