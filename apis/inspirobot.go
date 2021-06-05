package apis

import (
	"io/ioutil"
	"net/http"
)

const (
	generate_url = "https://inspirobot.me/api?generate=true"
)

func GetInspirobotMessage() (string, error) {
	resp, err := http.Get(generate_url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
