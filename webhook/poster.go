package webhook

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Poster struct {
	url    string
	client *http.Client
}

func NewPoster(url string) (poster Poster) {
	poster.url = url
	poster.client = &http.Client{
		Timeout: time.Second * 10,
	}
	return
}

func (poster Poster) Post(object interface{}) {
	jsonBytes, err := json.Marshal(object)
	if err != nil {
		panic(err)
	}
	jsonBuf := bytes.NewBuffer(jsonBytes)

	resp, err := poster.client.Post(poster.url, "application/json", jsonBuf)
	defer resp.Body.Close()
	log.Println(resp.Status)
}
