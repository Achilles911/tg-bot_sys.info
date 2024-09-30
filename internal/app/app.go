package app

import (
	"log"
	"os"
	"tgbot/internal/handlers"
	"time"

	tele "gopkg.in/telebot.v3"
)

func Run() {
	token := os.Getenv("TOKEN")
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	handlers.SetupHandlers(b)

	b.Start()
}
