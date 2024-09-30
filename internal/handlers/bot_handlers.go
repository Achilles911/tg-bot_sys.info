package handlers

import (
	"tgbot/internal/monitoring"
	"tgbot/internal/remote"

	tele "gopkg.in/telebot.v3"
)

func SetupHandlers(b *tele.Bot) {
	b.Handle("/remote_status", remote.HandleRemoteStatus)
	b.Handle("/status", monitoring.HandleStatus)
	b.Handle("/resources", monitoring.HandleResources)
}
