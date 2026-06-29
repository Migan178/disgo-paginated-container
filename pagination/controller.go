package pagination

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/rest"
)

type updatableEvent interface {
	discord.Interaction
	UpdateMessage(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error
	CreateFollowupMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*discord.Message, error)
}

func (p *PaginatedContainer) makeComponents() discord.ActionRowComponent {
	disabled := p.total == 1

	firstButton := discord.NewButton(
		p.opts.FirstButton.Style,
		p.opts.FirstButton.Label,
		makePaginationContainerFirst(p.id),
		"", 0,
	).
		WithEmoji(p.opts.FirstButton.Emoji).
		WithDisabled(disabled)

	prevButton := discord.NewButton(
		p.opts.PrevButton.Style,
		p.opts.PrevButton.Label,
		makePaginationContainerPrev(p.id),
		"", 0,
	).
		WithEmoji(p.opts.PrevButton.Emoji).
		WithDisabled(disabled)

	statusButton := discord.NewButton(
		p.opts.StatusButton.Style,
		fmt.Sprintf(p.opts.StatusButton.Label, p.current, p.total),
		makePaginationContainerPages(p.id),
		"", 0,
	).
		WithDisabled(!p.opts.StatusButton.EnableClickToSetPage || disabled)

	nextButton := discord.NewButton(
		p.opts.NextButton.Style,
		p.opts.NextButton.Label,
		makePaginationContainerNext(p.id),
		"", 0,
	).
		WithEmoji(p.opts.NextButton.Emoji).
		WithDisabled(disabled)

	lastButton := discord.NewButton(
		p.opts.LastButton.Style,
		p.opts.LastButton.Label,
		makePaginationContainerLast(p.id),
		"", 0,
	).
		WithEmoji(p.opts.LastButton.Emoji).
		WithDisabled(disabled)

	if p.opts.StatusButton.Emoji != nil {
		statusButton.Emoji = p.opts.StatusButton.Emoji
	}

	return discord.NewActionRow(firstButton, prevButton, statusButton, nextButton, lastButton)
}

func (p *PaginatedContainer) toFirst(i *handler.ComponentEvent) error {
	if p.current == 1 {
		return p.set(i, p.total)
	}

	return p.set(i, 1)
}

func (p *PaginatedContainer) toPrev(i *handler.ComponentEvent) error {
	if p.current == 1 {
		return p.set(i, p.total)
	}

	return p.set(i, p.current-1)
}

func (p *PaginatedContainer) toNext(i *handler.ComponentEvent) error {
	if p.current == p.total {
		return p.set(i, 1)
	}

	return p.set(i, p.current+1)
}

func (p *PaginatedContainer) toLast(i *handler.ComponentEvent) error {
	if p.current == p.total {
		return p.set(i, 1)
	}

	return p.set(i, p.total)
}

func (p *PaginatedContainer) set(i updatableEvent, page int) error {
	p.resetTimer()
	p.setPage(page)

	container := p.containers[p.current-1].AddComponents(p.makeComponents())
	return i.UpdateMessage(
		discord.NewMessageUpdateV2([]discord.LayoutComponent{container}...),
	)
}

func (p *PaginatedContainer) showModal(i *handler.ComponentEvent) error {
	return i.Modal(
		discord.NewModalCreate(
			makePaginationContainerModal(p.id),
			p.opts.Modal.Title,
			[]discord.LayoutComponent{
				discord.NewLabel(
					p.opts.Modal.Label,
					discord.NewShortTextInput(paginationContainerSetPage).
						WithPlaceholder(p.opts.Modal.Placeholder).
						WithValue(fmt.Sprint(p.current)),
				).
					WithDescription(p.opts.Modal.Description),
			}...,
		),
	)
}
