package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type APIResponse interface{}

// Respostas das APIs
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
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

func main() {
	// Criando canais para comunicação entre goroutines
	channelViaCEP := make(chan ViaCEP)
	channelApiCEP := make(chan ApiCEP)

	// Iniciando as goroutines para as requisições simultâneas
	go GetViaCEP(channelViaCEP)
	go GetApiCEP(channelApiCEP)

	// Seleção da primeira resposta que chegar
	select {
	case resViaCEP := <-channelViaCEP:
		// Exibe a resposta da ViaCEP
		fmt.Printf("ViaCEP: %+v\n", resViaCEP)
	case resApiCEP := <-channelApiCEP:
		// Exibe a resposta da ApiCEP
		fmt.Printf("ApiCEP: %+v\n", resApiCEP)
	case <-time.After(time.Second):
		// Caso não haja resposta dentro de 1 segundo
		fmt.Println("Timeout")
	}
}

// Função para buscar dados do ViaCEP
func GetViaCEP(chApi chan ViaCEP) {
	var viaCEP ViaCEP
	err := RequestAPI("https://viacep.com.br/ws/22621252/json/", &viaCEP)
	if err != nil {
		fmt.Println("Erro na API ViaCEP:", err)
		chApi <- ViaCEP{}
		return
	}
	viaCEP.ApiResponse.Api = "ViaCEP"
	chApi <- viaCEP
}

// Função para buscar dados do ApiCEP
func GetApiCEP(chVia chan ApiCEP) {
	var apiCEP ApiCEP
	err := RequestAPI("https://cdn.apicep.com/file/apicep/22260-003.json", &apiCEP)
	if err != nil {
		fmt.Println("Erro na API ApiCEP:", err)
		chVia <- ApiCEP{}
		return
	}
	apiCEP.ApiResponse.Api = "ApiCEP"
	chVia <- apiCEP
}

// Função genérica para requisições HTTP com controle de timeout
func RequestAPI(url string, res APIResponse) error {
	// Configurando cliente HTTP com timeout
	client := &http.Client{
		Timeout: 1 * time.Second, // Timeout de 1 segundo
	}

	// Realizando a requisição GET
	req, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("erro ao realizar a requisição: %w", err)
	}
	defer req.Body.Close()

	// Lendo o corpo da resposta
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("erro ao ler o corpo da resposta: %w", err)
	}

	// Deserializando a resposta JSON
	err = json.Unmarshal(body, res)
	if err != nil {
		return fmt.Errorf("erro ao decodificar a resposta JSON: %w", err)
	}

	return nil
}
