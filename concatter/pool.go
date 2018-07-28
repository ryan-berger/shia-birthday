package concatter

import (
	"net/http"
	"image/gif"
	"github.com/pborman/uuid"
	"os"
	"fmt"
	"encoding/json"
	"bytes"
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
	httpRequests map[string]string
}

func (pool *WorkerPool) MakeRequest(text string, responseUrl string) {
	id := uuid.NewUUID()
	pool.httpRequests[id.String()] = responseUrl
	pool.requests <- &gifRequest{
		id:   id,
		text: text,
	}
}

func (pool *WorkerPool) dispatchCalls() {
	for {
		select {
		case res := <-pool.results:
			responseUrl := pool.httpRequests[res.requestId.String()]
			f, e := os.Create(fmt.Sprintf("gifs/%s.gif", res.text))

			if e != nil {
				fmt.Println(e)
			}

			gif.EncodeAll(f, res.result)

			http.Post(responseUrl, "application/json", getSlackResponse(res.text))
		}
	}
}

func NewWorkerPool() *WorkerPool {
	pool := &WorkerPool{
		requests:     make(chan *gifRequest),
		results:      make(chan *gifResult),
		httpRequests: map[string]string{},
	}

	for i := 0; i < 2; i++ {
		pool.workers = append(pool.workers, NewWorker(pool.requests, pool.results))
	}

	go pool.dispatchCalls()
	return pool
}
