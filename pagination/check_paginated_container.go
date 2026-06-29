package pagination

import (
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func getID(customID string) string {
	parts := strings.Split(customID, "/")
	return parts[len(parts)-1]
}

func checkPaginationContainer() handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return func(e *handler.InteractionEvent) error {
			var customID string

			inter, ok := e.Interaction.(discord.ComponentInteraction)
			if !ok {
				inter, ok := e.Interaction.(discord.ModalSubmitInteraction)
				if !ok {
					return nil
				}

				customID = inter.Data.CustomID
			} else {
				customID = inter.Data.CustomID()
			}

			id := getID(customID)
			if getPaginationContainer(id) == nil {
				return nil
			}

			return next(e)
		}
	}
}
