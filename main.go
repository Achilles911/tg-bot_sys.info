package main

import (
	"log"
	"os/exec"
	"time"

	tele "gopkg.in/telebot.v3"
)

func main() {
	pref := tele.Settings{
		Token:  "7626218627:AAEz8Hrtr8MNtdcjt9iaog5zJsd5agDBVCI",
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Команда /status для проверки запущенных процессов
	b.Handle("/status", func(c tele.Context) error {
		// Выполняем команду ps aux
		out, err := exec.Command("ps", "aux").Output()
		if err != nil {
			return c.Send("Ошибка при выполнении команды ps aux")
		}

		// Конвертируем вывод в строку
		output := string(out)
		// Ограничиваем длину сообщения
		if len(output) > 4000 {
			output = output[:4000] + "\n...Output truncated..."
		}

		// Отправляем результат
		return c.Send("Запущенные процессы:\n" + output)
	})

	b.Start()
}
