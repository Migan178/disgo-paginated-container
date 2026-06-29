package pagination

import (
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func Router() handler.Router {
	r := handler.New()

	r.Use(checkPaginationContainer())

	// To first
	r.Component(paginationContainerFirst+"/{id}", func(e *handler.ComponentEvent) error {
		return getPaginationContainer(e.Vars["id"]).toFirst(e)
	})

	// To previous
	r.Component(paginationContainerPrev+"/{id}", func(e *handler.ComponentEvent) error {
		return getPaginationContainer(e.Vars["id"]).toPrev(e)
	})

	// To next
	r.Component(paginationContainerNext+"/{id}", func(e *handler.ComponentEvent) error {
		return getPaginationContainer(e.Vars["id"]).toNext(e)
	})

	// To last
	r.Component(paginationContainerLast+"/{id}", func(e *handler.ComponentEvent) error {
		return getPaginationContainer(e.Vars["id"]).toLast(e)
	})

	// Shows Modal
	r.Component(paginationContainerPages+"/{id}", func(e *handler.ComponentEvent) error {
		return getPaginationContainer(e.Vars["id"]).showModal(e)
	})

	// Set page
	r.Modal(paginationContainerModal+"/{id}", func(e *handler.ModalEvent) error {
		p := getPaginationContainer(e.Vars["id"])

		page, err := strconv.Atoi(e.Data.Text(paginationContainerSetPage))
		if err != nil {
			return e.CreateMessage(
				discord.NewMessageCreate().WithContent("해당 값은 숫자여야해요.").
					WithEphemeral(true),
			)
		}

		return p.set(e, page)
	})

	return r
}
