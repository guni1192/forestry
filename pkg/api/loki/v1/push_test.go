package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/guni1192/forestry/pkg/api"
	"github.com/stretchr/testify/assert"
)

func NewTestLogData() LogData {
	nowNS := strconv.Itoa(time.Now().Nanosecond())
	values := [][]string{{nowNS, "hoge", "foo", "bar"}}
	stream := Stream{Stream: map[string]string{"job": "test-server"}, Values: values}

	return LogData{Streams: []Stream{stream}}
}

func TestLokiV1PushHandlerShouldCreated(t *testing.T) {
	router := api.NewServer(false)
	router.POST("/api/loki/v1/push", PushHandler)

	logData := NewTestLogData()
	reqBody, err := json.Marshal(logData)
	if err != nil {
		assert.Fail(t, "Failed to encoding HTTP request body", err)
	}

	req := httptest.NewRequest("POST", "/api/loki/v1/push", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	body := make(map[string]string)
	err = json.Unmarshal(rec.Body.Bytes(), &body)
	if err != nil {
		assert.Fail(t, "Failed to encoding HTTP response body", err)
	}

	assert.Equal(t, "Pushing log data succeseed", body["message"])
}

func TestLokiV1PushHandlerShouldBadRequest(t *testing.T) {
	router := api.NewServer(false)
	router.POST("/api/loki/v1/push", PushHandler)

	dummy, err := json.Marshal(map[string]string{"dummy": "data"})
	if err != nil {
		assert.Failf(t, "Failed to encoding dummy request body: %s", err.Error())
	}

	req := httptest.NewRequest("POST", "/api/loki/v1/push", bytes.NewBuffer(dummy))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	body := make(map[string]interface{})
	err = json.Unmarshal(rec.Body.Bytes(), &body)
	if err != nil {
		assert.Failf(t, "Failed to encoding HTTP response body: %s", err.Error())
	}

	assert.Equal(t, "Bad Request. logData is empty.", body["message"])
}

// Unset Content-Type: application/json
func TestLokiV1PushHandlerShouldUnsupportedMediaType(t *testing.T) {
	router := api.NewServer(false)
	router.POST("/api/loki/v1/push", PushHandler)

	dummy, err := json.Marshal(map[string]string{"dummy": "data"})
	if err != nil {
		assert.Failf(t, "Failed to encoding dummy request body: %s", err.Error())
	}

	req := httptest.NewRequest("POST", "/api/loki/v1/push", bytes.NewBuffer(dummy))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)

	body := make(map[string]interface{})
	err = json.Unmarshal(rec.Body.Bytes(), &body)
	if err != nil {
		assert.Failf(t, "Failed to encoding HTTP response body: %s", err.Error())
	}

	assert.Equal(t, "code=415, message=Unsupported Media Type", body["message"])

}
