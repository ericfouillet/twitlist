package twitlist

import "github.com/eric-fouillet/twitlist/twitlistserver"

func init() {
	twitlistserver.RegisterHandlers()
}
