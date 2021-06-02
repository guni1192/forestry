package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	lokiV1 "github.com/guni1192/forestry/pkg/api/loki/v1"
	"github.com/spf13/cobra"
)

var (
	rootCommand = &cobra.Command{
		Use:   "forestry-ag",
		Short: "forestry agent program",
		RunE: func(cmd *cobra.Command, args []string) error {
			forestryHost, err := cmd.Flags().GetString("forestry-host")
			logPath, err := cmd.Flags().GetString("log-file")
			if err != nil {
				return err
			}

			if logPath == "" {
				return errors.New("Please specify --log-file")
			}

			run(forestryHost, logPath)

			return nil
		},
	}
)

type forestryClient struct {
	forestryHost string
}

func init() {
	rootCommand.PersistentFlags().StringP("forestry-host", "", "http://localhost:1192", "Forestry hostname")
	rootCommand.PersistentFlags().StringP("log-file", "f", "", "target log file")
}

func (c *forestryClient) send(message string) error {
	now := int(time.Now().Unix() * 1000000000)
	t := strconv.Itoa(now)

	values := make([][]string, 1)
	values[0] = []string{t, message}

	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("Failed to get hostname: %w", err)
	}

	stream := lokiV1.Stream{
		Stream: map[string]string{"hostname": hostname},
		Values: values,
	}

	pushData := lokiV1.LogData{
		Streams: []lokiV1.Stream{stream},
	}

	data, err := json.Marshal(pushData)
	if err != nil {
		return fmt.Errorf("Failed to marshal json data: %w", err)
	}

	url := fmt.Sprintf("%s/api/loki/v1/push", c.forestryHost)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("POST %s failed: %w", url, err)
	}
	defer resp.Body.Close()

	return nil
}

func (c *forestryClient) tail(file io.Reader) error {
	r := bufio.NewReader(file)

	for {
		bytes, err := r.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return fmt.Errorf("Failed to read buffer: %w", err)
		}
		if len(bytes) != 0 {
			if err != nil {
				c.send(string(bytes))
			}
		}
		if err == io.EOF {
			time.Sleep(time.Millisecond * 50)
		}
	}
}

func run(forestryHost string, logPath string) error {
	f, err := os.Open(logPath)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("Failed to open log data %w", err)
	}

	c := forestryClient{forestryHost: forestryHost}

	c.tail(f)
	if err != nil {
		return fmt.Errorf("Failed to tail log data %w", err)
	}

	return nil
}

func main() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
