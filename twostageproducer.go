package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type TwoStageAnswerProducer struct {
	RetrieverUrl  string `json:"retriever-url"`
	SummarizerUrl string `json:"summarizer-url"`
}

func NewTwoStageAnswerProducer(config string) (AnswerProducer, error) {
	ap := &TwoStageAnswerProducer{}

	byt, err := ioutil.ReadFile(config)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(byt, ap); err != nil {
		return nil, err
	}

	log.Printf("twostageproducer: RetrieverUrl `%s`\n", ap.RetrieverUrl)
	log.Printf("twostageproducer: SummarizerUrl `%s`\n", ap.SummarizerUrl)
	return ap, nil
}

func (ap *TwoStageAnswerProducer) GetAnswer(result chan *Answer, q *Question) {
	var answer *Answer
	var summary string
	var resources []string
	var passages []string

	retrievers := []Retriever{
		NewRemoteRetriever(ap.RetrieverUrl),
		NewDummyRetriever(),
	}

	summarizers := []Summarizer{
		NewRemoteSummarizer(ap.SummarizerUrl),
		NewDummySummarizer(),
	}

	var err error

	for _, retriever := range retrievers {
		passages, err = retriever.GetPassages(q, 10)
		if err != nil {
			answer = NewErrorAnswer(q, err)
			continue
		}
		break
	}

	for _, summarizer := range summarizers {
		summary, err = summarizer.GetSummary(passages, q, config.AnswerSize)
		if err != nil {
			answer = NewErrorAnswer(q, err)
			continue
		}

		answer = &Answer{
			Answered:  "yes",
			Pid:       config.Pid,
			Qid:       q.Qid,
			Time:      int64(time.Since(q.ReceivedTime) / time.Millisecond),
			Content:   Truncate(summary, config.AnswerSize),
			Resources: strings.Join(resources, ","),
		}
		goto end
	}

end:
	result <- answer
}
