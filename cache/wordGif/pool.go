package wordGif

import (
			"github.com/pborman/uuid"
		"fmt"
	"encoding/json"
	"bytes"
		"image/gif"
	)

type attachment struct {
	Text     string `json:"text"`
	ImageUrl string `json:"image_url"`
}

type slackResponse struct {
	ResponseType string       `json:"response_type"`
	Attachments  []attachment `json:"attachments"`
}

func getSlackResponse(text string) *bytes.Reader {
	resp := slackResponse{
		ResponseType: "in_channel",
		Attachments: []attachment{
			{Text: text, ImageUrl: fmt.Sprintf("http://slack.ryanberger.me/gif/%s", text)},
		},
	}
	response, _ := json.Marshal(&resp)

	return bytes.NewReader(response)
}

type WorkerPool struct {
	requests     chan *gifRequest
	results      chan *gifResult
	workers      []*Worker
	resultChan map[string]chan *gif.GIF
}

func (pool *WorkerPool) MakeRequest(text string, resultChan chan *gif.GIF) {
	id := uuid.NewUUID()
	pool.resultChan[id.String()] = resultChan
	pool.requests <- &gifRequest{
		id:   id,
		text: text,
	}
}

func (pool *WorkerPool) dispatchCalls() {
	for {
		select {
		case res := <-pool.results:
			pool.resultChan[res.requestId.String()] <- res.result
		}
	}
}

func NewWorkerPool() *WorkerPool {
	pool := &WorkerPool{
		requests:     make(chan *gifRequest),
		results:      make(chan *gifResult),
		resultChan : map[string] chan *gif.GIF {},
	}

	for i := 0; i < 2; i++ {
		pool.workers = append(pool.workers, NewWorker(pool.requests, pool.results))
	}

	//go pool.dispatchCalls()
	return pool
}
