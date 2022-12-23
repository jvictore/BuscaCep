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

func main () {
	http.HandleFunc("/", BuscaCepHandler)
	http.ListenAndServe(":8080", nil)
}

func BuscaCepHandler(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/"{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	cepParam := r.URL.Query().Get("cep")
	if cepParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cep, error := BuscaCep(cepParam)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(*cep)

}

func BuscaCep(cep string) (*ViaCEP, error) { 
	var urlLeft string = "http://viacep.com.br/ws/"
	var urlRight string = "/json/"
	
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

	return &dataCep, nil
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