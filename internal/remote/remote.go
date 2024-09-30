package remote

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	tele "gopkg.in/telebot.v3"
)

type ServiceResponse struct {
	Status      string `json:"status"`
	Message     string `json:"message"`
	CPUUsage    string `json:"cpu_usage"`
	MemoryUsage string `json:"memory_usage"`
}

func HandleRemoteStatus(c tele.Context) error {
	serviceURL := os.Getenv("SERVICEURL")

	resp, err := http.Get(serviceURL)
	if err != nil {
		return c.Send("Ошибка при подключении к удалённому сервису")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Send("Ошибка при чтении ответа от сервиса")
	}

	var serviceResponse ServiceResponse
	err = json.Unmarshal(body, &serviceResponse)
	if err != nil {
		return c.Send("Ошибка при парсинге данных от сервера")
	}

	output := fmt.Sprintf("== Статус удалённого сервера ==\nService: %s\nCPU: %s\nMemory: %s\nStatus: %s",
		serviceResponse.Message, serviceResponse.CPUUsage, serviceResponse.MemoryUsage, serviceResponse.Status)

	return c.Send(output)
}
