package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Migan178/disgo-paginated-container/pagination"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func main() {
	session, err := disgo.New("INPUT TOKEN HERE",
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
				gateway.IntentGuildMessages,
				gateway.IntentMessageContent,
			),
		),
		bot.WithEventListeners(pagination.Router()),
		bot.WithEventListenerFunc(onMessageCreate),
	)
	if err != nil {
		panic(err)
	}

	if err := session.OpenGateway(context.Background()); err != nil {
		panic(err)
	}

	defer session.Close(context.Background())

	fmt.Println("bot is running. press ctrl+C to exit program.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func onMessageCreate(e *events.MessageCreate) {
	if e.Message.Author.Bot {
		return
	}

	if e.Message.Content == "ping" {
		if _, err := e.Client().Rest.CreateMessage(e.ChannelID, discord.NewMessageCreate().WithContent("Pong!")); err != nil {
			fmt.Println(err)
		}

		return
	}

	if e.Message.Content == "page" {
		if err := pagination.New(e, false, nil,
			discord.NewContainer(discord.NewTextDisplay("# first")),
			discord.NewContainer(discord.NewTextDisplay("# second")),
		).Start(); err != nil {
			fmt.Println(err)
		}

		return
	}
}
