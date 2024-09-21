package main

import (
	"log"
	"os/exec"
	"strings"
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

	monitoredProcesses := []string{"nginx", "mysql", "./bot"} //Cписок процессов, по которым будет фильрация

	// Команда /status для проверки запущенных процессов
	b.Handle("/status", func(c tele.Context) error {
		// Выполняем команду ps aux
		out, err := exec.Command("ps", "aux").Output()
		if err != nil {
			return c.Send("Ошибка при выполнении команды ps aux")
		}
		output := string(out)
		lines := strings.Split(output, "\n")
		// Фильтр строки по отслеживаемым процессам
		var filteredLines []string
		for _, line := range lines {
			for _, process := range monitoredProcesses {
				if strings.Contains(line, process) {
					filteredLines = append(filteredLines, line)
					break
				}
			}
		}

		if len(filteredLines) == 0 {
			return c.Send("Нет запущенных процессов из отслеживаемого списка.")
		}

		//Процессы в строку
		output = strings.Join(filteredLines, "\n")

		// Ограничиваем длину сообщения
		if len(output) > 4000 {
			output = output[:4000] + "\n...Output truncated..."
		}

		// Отправляем результат
		return c.Send("Запущенные процессы:\n" + output)
	})

	b.Start()
}
