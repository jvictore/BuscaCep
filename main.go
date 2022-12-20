package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ViaCEP struct {
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

func main() {
	var ceps [] string
	numParams := len(os.Args) - 1

	// Populate the ceps array
	if numParams == 0 {
		var numCeps int
		var cep string
		
		println("How many CEPs you will check?")
		fmt.Scan(&numCeps)

		for i := 0; i < numCeps; i++{
			fmt.Print("Enter the ", i+1, " CEP: ")
			fmt.Scan(&cep)
			ceps = append(ceps, cep)
		}

	} else {
		for _, cep := range os.Args[1:] {
			ceps = append(ceps, cep)
		}
	}

	var urlLeft string = "http://viacep.com.br/ws/"
	var urlRight string = "/json/"

	for idxCep, cep := range ceps {
		req, err := http.Get(urlLeft + cep + urlRight)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error executing the request: %v\n", err)
		}
		defer req.Body.Close()
		
		res, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading the response: %v\n", err)
		}

		var dataCep ViaCEP
		err = json.Unmarshal(res, &dataCep)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing the response: %v\n", err)
		}

		printDataCep(idxCep, &dataCep)
		
	}

	
}

func printDataCep(idxCep int, dataCep *ViaCEP){
	fmt.Println(idxCep+1, "CEP:")
	fmt.Println("CEP: ", dataCep.Cep)
	fmt.Println("Logradouro: ", dataCep.Logradouro)
	fmt.Println("Complemento: ", dataCep.Complemento)
	fmt.Println("Bairro: ", dataCep.Bairro)
	fmt.Println("Localidade: ", dataCep.Localidade)
	fmt.Println("UF: ", dataCep.Uf)
	fmt.Println("IBGE: ", dataCep.Ibge)
	fmt.Println()
}