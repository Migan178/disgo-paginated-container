package pagination

import (
	"time"

	"github.com/disgoorg/disgo/discord"
)

// PaginatedContainerOptions is a options of paginated container.
type PaginatedContainerOptions struct {
	EndDuration time.Duration

	FirstButton  *ButtonStyle
	PrevButton   *ButtonStyle
	StatusButton *StatusButtonStyle
	NextButton   *ButtonStyle
	LastButton   *ButtonStyle

	Modal *ModalOption

	PageSetErrorMessage string
}

// ButtonStyle is how to show button.
type ButtonStyle struct {
	Label string
	Emoji discord.ComponentEmoji
	Style discord.ButtonStyle
}

// StatusButtonStyle is how to show status button.
type StatusButtonStyle struct {
	// It should have 2 of %d(First %d is current page, Last %d is total pages).
	Label                string
	Emoji                *discord.ComponentEmoji
	Style                discord.ButtonStyle
	EnableClickToSetPage bool
}

// ModalOption is how to show modal.
type ModalOption struct {
	Title       string
	Label       string
	Placeholder string
	Description string
}

func defaultOpts() *PaginatedContainerOptions {
	return &PaginatedContainerOptions{
		EndDuration: 10 * time.Minute,
		FirstButton: &ButtonStyle{
			Emoji: discord.NewComponentEmoji("⏪"),
			Style: discord.ButtonStylePrimary,
		},
		PrevButton: &ButtonStyle{
			Emoji: discord.NewComponentEmoji("◀️"),
			Style: discord.ButtonStylePrimary,
		},
		StatusButton: &StatusButtonStyle{
			Label:                "(%d/%d)",
			Style:                discord.ButtonStyleSecondary,
			EnableClickToSetPage: true,
		},
		NextButton: &ButtonStyle{
			Emoji: discord.NewComponentEmoji("▶️"),
			Style: discord.ButtonStylePrimary,
		},
		LastButton: &ButtonStyle{
			Emoji: discord.NewComponentEmoji("⏩"),
			Style: discord.ButtonStylePrimary,
		},

		Modal: &ModalOption{
			Title:       "Set Page",
			Label:       "Page to move",
			Placeholder: "input the page number here...",
			Description: "Input the page number please.",
		},

		PageSetErrorMessage: "The value must be an number!",
	}
}
