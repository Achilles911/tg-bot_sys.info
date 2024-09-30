package monitoring

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v3"
)

var monitoredProcesses = []string{"nginx", "mysql", "./bot"}

func HandleStatus(c tele.Context) error {
	out, err := exec.Command("top", "-b", "-n", "1").Output()
	if err != nil {
		return c.Send("Ошибка при выполнении команды top")
	}
	output := string(out)
	lines := strings.Split(output, "\n")

	header := fmt.Sprintf("%-3s %-2s %-5s %-4s %-5s %-1s %-1s %-0s %-s\n", "PID", "PPID", "USER", "STAT", "VSZ", "%VSZ", "CPU", "%CPU", "COMMAND")

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
		fields := strings.Fields(line)
		if len(fields) >= 10 {
			formattedLine := fmt.Sprintf("%-6s %-8s %-8s %-6s %-8s %-5s %-4s %-4s %-s",
				fields[0], fields[1], fields[2], fields[3], fields[4], fields[5], fields[6], fields[7], strings.Join(fields[10:], " "))
			formattedOutput = append(formattedOutput, formattedLine)
		}
	}

	output = header + strings.Join(formattedOutput, "\n")

	if len(output) > 4000 {
		output = output[:4000] + "\n...Output truncated..."
	}

	return c.Send("Запущенные процессы:\n" + output)
}

func HandleResources(c tele.Context) error {
	topOut, err := exec.Command("top", "-b", "-n", "1").Output()
	if err != nil {
		return c.Send("Ошибка при получении данных о CPU")
	}

	freeOut, err := exec.Command("free", "-b").Output()
	if err != nil {
		return c.Send("Ошибка при получении данных о памяти")
	}

	lines := strings.Split(string(topOut), "\n")
	var cpuLine string
	for _, line := range lines {
		if strings.Contains(line, "CPU:") {
			cpuLine = line
			break
		}
	}

	freeLines := strings.Split(string(freeOut), "\n")
	var memLine string
	for _, line := range freeLines {
		if strings.Contains(line, "Mem:") {
			memLine = line
			break
		}
	}

	memFields := strings.Fields(memLine)
	if len(memFields) >= 7 {
		usedMem, _ := strconv.ParseFloat(memFields[1], 64)
		freeMem, _ := strconv.ParseFloat(memFields[2], 64)
		buffMem, _ := strconv.ParseFloat(memFields[4], 64)
		cachedMem, _ := strconv.ParseFloat(memFields[5], 64)

		usedMemMB := usedMem / (1024 * 1024)
		freeMemMB := freeMem / (1024 * 1024)
		buffMemMB := buffMem / (1024 * 1024)
		cachedMemMB := cachedMem / (1024 * 1024)

		memLine = fmt.Sprintf("Mem: %.2f MB used, %.2f MB free, %.2f MB buff, %.2f MB cached",
			usedMemMB, freeMemMB, buffMemMB, cachedMemMB)
	}

	finalOutput := "Мониторинг ресурсов\n\n" + cpuLine + "\n" + memLine

	return c.Send(finalOutput)
}
