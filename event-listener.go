package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type triggerSyncBody struct {
	JobType      string `json:"jobType"`
	ConnectionID string `json:"connectionId"`
}

func customersUpdated(airbyteURL, airbyteConnectionID, basicAuth string) func(ctx context.Context, msg []byte) error {
	return func(ctx context.Context, msg []byte) error {
		tsb := triggerSyncBody{
			JobType:      "sync",
			ConnectionID: airbyteConnectionID,
		}

		payload, _ := json.Marshal(tsb)
		bb := bytes.NewBuffer(payload)
		req, err := http.NewRequest(http.MethodPost, airbyteURL, bb)
		if err != nil {
			return err
		}

		req.Header.Add("accept", "application/json")
		req.Header.Add("content-type", "application/json")
		req.Header.Add("authorization", basicAuth)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		log.Printf("sync request at connection %s status: %s\n", airbyteConnectionID, res.Status)
		return nil
	}
}
