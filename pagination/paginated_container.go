package pagination

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"
)

// PaginatedContainer is container with page
type PaginatedContainer struct {
	containers []discord.ContainerComponent
	current    int
	total      int
	id         string
	timer      *time.Timer
	deferred   bool
	e          any
	opts       *PaginatedContainerOptions
	messageID  snowflake.ID
}

type creatableInteractionEvent interface {
	discord.Interaction
	Client() *bot.Client
	CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error
}

var paginatedContainers = make(map[string]*PaginatedContainer)

// New creates a new PaginationContainer
func New(event any, deferred bool, opts *PaginatedContainerOptions, containers ...discord.ContainerComponent) *PaginatedContainer {
	var userID snowflake.ID

	if msgEvent, ok := event.(*events.MessageCreate); ok {
		userID = msgEvent.Message.Author.ID
	}

	if interEvent, ok := event.(creatableInteractionEvent); ok {
		userID = interEvent.User().ID
	}

	if opts == nil {
		opts = defaultOpts()
	}

	id := fmt.Sprintf("%d:%s", userID, rand.Text())
	return &PaginatedContainer{
		containers: containers,
		current:    1,
		total:      len(containers),
		id:         id,
		timer:      time.NewTimer(opts.EndDuration),
		deferred:   deferred,
		e:          event,
		opts:       opts,
	}
}

func (p *PaginatedContainer) waitTimerEnd() {
	<-p.timer.C

	msgUpdate := discord.NewMessageUpdateV2([]discord.LayoutComponent{
		p.containers[p.current-1],
	}...)

	if msgEvent, ok := p.e.(*events.MessageCreate); ok {
		msgEvent.Client().Rest.UpdateMessage(msgEvent.ChannelID, p.messageID, msgUpdate)
	}

	if interEvent, ok := p.e.(creatableInteractionEvent); ok {
		interEvent.Client().Rest.UpdateInteractionResponse(
			interEvent.ApplicationID(),
			interEvent.Token(),
			msgUpdate,
		)
	}

	delete(paginatedContainers, p.id)
}

func (p *PaginatedContainer) resetTimer() {
	p.timer.Reset(p.opts.EndDuration)
}

func (p *PaginatedContainer) setPage(page int) {
	if page <= 0 {
		p.current = 1
	} else if page > p.total {
		p.current = p.total
	} else {
		p.current = page
	}
}

// AddContainers adds containers.
func (p *PaginatedContainer) AddContainers(containers ...discord.ContainerComponent) *PaginatedContainer {
	p.total += len(containers)
	p.containers = append(p.containers, containers...)
	return p
}

// SetStartPage sets start page number.
func (p *PaginatedContainer) SetStartPage(startPage int) *PaginatedContainer {
	p.setPage(startPage)
	return p
}

// Start starts the paginated container.
func (p *PaginatedContainer) Start() error {
	if len(p.containers) == 0 {
		return nil
	}

	container := p.containers[p.current-1].AddComponents(p.makeComponents())
	paginatedContainers[p.id] = p

	go p.waitTimerEnd()

	if msgEvent, ok := p.e.(*events.MessageCreate); ok {
		msg, err := msgEvent.Client().Rest.CreateMessage(msgEvent.ChannelID, discord.NewMessageCreateV2(container))
		if err != nil {
			return err
		}

		p.messageID = msg.ID
		return nil
	}

	if interEvent, ok := p.e.(creatableInteractionEvent); ok {
		if p.deferred {
			_, err := interEvent.Client().Rest.UpdateInteractionResponse(
				interEvent.ApplicationID(),
				interEvent.Token(),
				discord.NewMessageUpdateV2([]discord.LayoutComponent{container}...),
			)
			return err
		}

		return interEvent.CreateMessage(discord.NewMessageCreateV2(container).WithEphemeral(true))
	}

	return nil
}

func getPaginationContainer(id string) *PaginatedContainer {
	if p, ok := paginatedContainers[id]; ok {
		return p
	}
	return nil
}
