package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

func main() {
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

	monitoredProcesses := []string{"nginx", "mysql", "./bot"} //Cписок процессов, по которым будет фильтрация

	// Команда /status для проверки запущенных процессов
	b.Handle("/status", func(c tele.Context) error {
		// Выполняем команду ps aux
		out, err := exec.Command("top", "-b", "-n", "1").Output()
		if err != nil {
			return c.Send("Ошибка при выполнении команды ps aux")
		}
		output := string(out)
		lines := strings.Split(output, "\n")

		header := fmt.Sprintf("%-3s %-2s %-5s %-4s %-5s %-1s %-1s %-0s %-s\n", "PID", "PPID", "USER", "STAT", "VSZ", "%VSZ", "CPU", "%CPU", "COMMAND")

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
		var formattedOutput []string
		for _, line := range filteredLines {
			fields := strings.Fields(line) // Разделяем строку на части (столбцы)
			if len(fields) >= 10 {         // Убедимся, что есть минимум 10 столбцов
				// Форматируем строку по размерам столбцов
				formattedLine := fmt.Sprintf("%-6s %-8s %-8s %-6s %-8s %-5s %-4s %-4s %-s",
					fields[0],                      // PID
					fields[1],                      // PPID
					fields[2],                      // USER
					fields[3],                      // STAT (вместо TTY)
					fields[4],                      // VSZ
					fields[5],                      // %VSZ
					fields[6],                      // CPU
					fields[7],                      // %CPU
					strings.Join(fields[10:], " ")) // COMMAND
				formattedOutput = append(formattedOutput, formattedLine)
			}
		}

		//Процессы в строку
		output = header + strings.Join(formattedOutput, "\n")

		// Ограничиваем длину сообщения
		if len(output) > 4000 {
			output = output[:4000] + "\n...Output truncated..."
		}

		// Отправляем результат
		return c.Send("Запущенные процессы:\n" + output)
	})

	b.Start()
}
