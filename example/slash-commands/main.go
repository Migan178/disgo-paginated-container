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
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
)

func router() handler.Router {
	r := handler.New()

	r.SlashCommand("/ping", handlePingCommand)
	r.SlashCommand("/page", handlePageCommand)

	return r
}

func main() {
	session, err := disgo.New("INPUT TOKEN HERE",
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
			),
		),
		bot.WithEventListeners(pagination.Router()),
		bot.WithEventListeners(router()),
	)
	if err != nil {
		panic(err)
	}

	if err := session.OpenGateway(context.Background()); err != nil {
		panic(err)
	}

	defer session.Close(context.Background())

	if _, err := session.Rest.SetGlobalCommands(session.ApplicationID, []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "ping",
			Description: "Pong!",
		},
		discord.SlashCommandCreate{
			Name:        "page",
			Description: "paginated container demo",
		},
	}); err != nil {
		panic(err)
	}

	fmt.Println("bot is running. press ctrl+C to exit program.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func handlePingCommand(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return e.CreateMessage(discord.NewMessageCreate().WithContent("Pong!"))
}

func handlePageCommand(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return pagination.New(e, false, nil,
		discord.NewContainer(discord.NewTextDisplay("# first")),
		discord.NewContainer(discord.NewTextDisplay("# second")),
	).Start()
}
