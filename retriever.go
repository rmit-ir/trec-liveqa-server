package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Retriever implements a simplistic interface for retrieval models.
// It takes a Question and an interger as input, the latter representing the
// maximum number of passages to retrieve.
type Retriever interface {
	GetPassages(q *Question, limit int) ([]string, error)
}

// DummyRetriever returns some canned responses as the passage
type DummyRetriever struct{}

func NewDummyRetriever() Retriever {
	return &DummyRetriever{}
}

func (dr *DummyRetriever) GetPassages(q *Question, limit int) ([]string, error) {
	resp := make([]string, 1)
	resp[0] = "That's your question, eh?"
	return resp, nil
}

// RemoteRetriever talks to a remote retrieval server in JSON
type RemoteRetriever struct{ Url string }

func NewRemoteRetriever(url string) Retriever {
	return &RemoteRetriever{Url: url}
}

// RemoteRetrieverRequest simply wraps around the Retriever protocol
type RemoteRetrieverRequest struct {
	Question Question `json:"quetion"`
	Limit    int      `json:"limit"`
}

type RemoteRetrieverResponse struct {
	Passages []string `json:"passages"`
}

func (rr *RemoteRetriever) GetPassages(q *Question, limit int) ([]string, error) {
	var passages []string

	req := &RemoteRetrieverRequest{*q, limit}
	payload, err := json.Marshal(req)
	if err != nil {
		return passages, err
	}

	rpayload, err := http.Post(rr.Url, "application/json", bytes.NewReader(payload))
	if err != nil {
		return passages, err
	}
	defer rpayload.Body.Close()

	body, err := ioutil.ReadAll(rpayload.Body)
	if err != nil {
		return passages, err
	}

	var resp RemoteRetrieverResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return passages, err
	}
	return resp.Passages, nil
}
