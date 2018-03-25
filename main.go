package main

import (
	"dwsb/actions"
	"dwsb/client"
	"log"
)

func main() {

	// Start Discord Client
	//client.Discord.Session.AddHandler(client.ReadyHandler)
	//client.Discord.Session.AddHandler(client.ChannelStateHandler)
	//_ = client.Discord.Session.Open()
	go client.Discord.Connect()

	app := actions.App()
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}

}
