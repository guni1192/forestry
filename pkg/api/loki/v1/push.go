package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Stream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

type LogData struct {
	Streams []Stream `json:"streams"`
}

type LogStreamer struct {
	lokiClient LokiClient
}

type LokiClient interface {
	Push(logData *LogData) error
}

type ForestryLokiClient struct {
	LokiHost string
}

func (c ForestryLokiClient) Push(logData *LogData) error {
	data, err := json.Marshal(logData)
	if err != nil {
		return fmt.Errorf("Failed to marshal json data: %w", err)
	}

	url := fmt.Sprintf("%s/loki/api/v1/push", c.LokiHost)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("POST %s failed: %w", url, err)
	}
	defer resp.Body.Close()

	return nil
}

func NewLogStreamer(lokiClient LokiClient) *LogStreamer {
	return &LogStreamer{lokiClient: lokiClient}
}

func (l *LogStreamer) PushHandler(c echo.Context) error {

	logData := new(LogData)

	err := c.Bind(&logData)
	if err != nil {
		return c.JSON(http.StatusUnsupportedMediaType, map[string]string{"message": err.Error()})
	}

	if len(logData.Streams) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Bad Request. logData is empty."})
	}

	err = l.lokiClient.Push(logData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Sorry. Internal Server Error"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Pushing log data succeseed"})
}
