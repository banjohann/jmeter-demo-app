package main

import (
	"github.com/JohannBandelow/jmeter-demo-app/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {

	viewsEngine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: viewsEngine,
	})

	app.Static("/static/", "./static")

	appHandler := handlers.NewAppHandler()
	wsServer := NewWebSocket()

	app.Get("/", appHandler.HandleGetIndex)
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	app.Get("/ws", websocket.New(func(ctx *websocket.Conn) {
		wsServer.HandleWebSocket(ctx)
	}))

	app.Post("/message", func(c *fiber.Ctx) error {
		msg := new(Message)

		if err := c.BodyParser(msg); err != nil {
			return err
		}

		wsServer.PublishMessage(msg)

		return c.Send([]byte(msg.ClientName))
	})

	go wsServer.HandleMessages()

	app.Listen(":3000")
}
