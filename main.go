package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type APIResponse interface{}

type ApiResponse struct {
	Api string `json:"Api"`
}

type ViaCEP struct {
	ApiResponse
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type ApiCEP struct {
	ApiResponse
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
	Status       int    `json:"status"`
	Ok           bool   `json:"ok"`
	StatusText   string `json:"statusText"`
}

func main() {
	channelViaCEP := make(chan ViaCEP)
	channelApiCEP := make(chan ApiCEP)

	go GetViaCEP(channelViaCEP)
	go GetApiCEP(channelApiCEP)

	select {
	case resViaCEP := <-channelViaCEP:
		fmt.Printf("ViaCEP: %+v\n", resViaCEP)

	case resApiCEP := <-channelApiCEP:
		fmt.Printf("Brasilapi: %+v\n", resApiCEP)

	case <-time.After(time.Second * 10):
		fmt.Printf("TimeOut\n")
	}
}

func GetViaCEP(chApi chan ViaCEP) {
	var viaCEP ViaCEP
	RequestAPI("https://viacep.com.br/ws/22621252/json/", &viaCEP)
	// Simular o tempo maior na requisição.
	// time.Sleep(2 * time.Second)
	viaCEP.ApiResponse.Api = "ViaCEP"
	chApi <- viaCEP
}

func GetApiCEP(chApi chan ApiCEP) {
	var apiCEP ApiCEP
	RequestAPI("https://brasilapi.com.br/api/cep/v1/08141140", &apiCEP)
	// Simular o tempo maior na requisição.
	// time.Sleep(2 * time.Second)
	apiCEP.ApiResponse.Api = "ApiCEP"
	chApi <- apiCEP
}

func RequestAPI(url string, res APIResponse) error {
	req, err := http.Get(url)
	if err != nil {
		return err
	}
	defer req.Body.Close()

	// Lê o corpo da resposta
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	// Verifica o código de status para erros como 429
	if req.StatusCode == 429 {
		return fmt.Errorf("erro 429: Too Many Requests")
	}

	// Deserializa o corpo JSON na estrutura fornecida
	err = json.Unmarshal(body, res)
	if err != nil {
		return err
	}
	return nil
}
