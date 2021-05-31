package v1

import (
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

func PushHandler(c echo.Context) error {
	logData := new(LogData)

	err := c.Bind(&logData)
	if err != nil {
		return c.JSON(http.StatusUnsupportedMediaType, map[string]string{"message": err.Error()})
	}

	if len(logData.Streams) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Bad Request. logData is empty."})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Pushing log data succeseed"})
}
