package service

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

type CepData struct {
	Provider string `json:"provider"`
	Response string `json:"response"`
}

func getApiData(url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("error on http.NewRequestWithContext(): %v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error on http.NewRequestWithContext(): %v", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error on io.ReadAll: %v", err)
	}

	return string(body), nil
}

func apicep(cep string, ch chan<- *CepData) {
	url := "https://cdn.apicep.com/file/apicep/" + cep + ".json"

	r, err := getApiData(url)
	if err != nil {
		return
	}

	d := CepData{
		Provider: "apicep",
		Response: r,
	}

	ch <- &d
}

func viacep(cep string, ch chan<- *CepData) {
	url := "https://viacep.com.br/ws/" + cep + "/json/"

	r, err := getApiData(url)
	if err != nil {
		return
	}

	d := CepData{
		Provider: "viacep",
		Response: r,
	}

	ch <- &d
}

func GetAddress(cep string) (*CepData, error) {
	validCep, _ := regexp.Compile(`^\d{5}-\d{3}$`)
	if !validCep.MatchString(cep) {
		return nil, fmt.Errorf("error: cep not in '99999-999' format")
	}

	ch := make(chan *CepData)

	go apicep(cep, ch)
	go viacep(cep, ch)

	select {
	case resp := <-ch:
		return resp, nil
	case <-time.After(1 * time.Second):
		return nil, fmt.Errorf("timeout for concurrent services")
	}
}
