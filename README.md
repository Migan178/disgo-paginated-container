# disgo-paginated-container

![Demo](assets/example.gif)

- A library that makes it easy to create paginated Containers using ComponentsV2 in [disgo](https://github.com/disgoorg/disgo).

# Installation

```sh
go get github.com/Migan178/disgo-paginated-container
```

# Usage

```go
package main

import (
    "github.com/Migan178/disgo-paginated-container/pagination"
    "github.com/disgoorg/disgo/discord"
    "github.com/disgoorg/disgo/events"
)

func onMessageCreate(e *events.MessageCreate) {
    if e.Message.Content == "page" {
        pagination.New(e,
            false, // deferred
            nil,   // options
            // containers
            discord.NewContainer(discord.NewTextDisplay("# first")),
            discord.NewContainer(discord.NewTextDisplay("# second")),
        ).Start()
        return
    }
}
```

- You can find more in the [examples](examples/) folder.
