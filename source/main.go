package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"github.com/google/uuid"
)

type ViaCEP struct {
	ID string
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func NewProduct(cep string, logradouro string, bairro string, localidade string, uf string, ibge string, gia, string, ddd string, siafi string) *ViaCEP{
	return &ViaCEP{
		ID: uuid.New().String(),
		Cep: cep,
		Logradouro: logradouro,
		Bairro: bairro,
		Localidade: localidade,
		Uf: uf,
		Ibge: ibge,
		Gia: gia,
		Ddd: ddd,
		Siafi: siafi,
	}
}

func main() {
	http.HandleFunc("/", SearchCepHandler)
	http.HandleFunc("/ui", SearchCepUI)
	http.ListenAndServe(":8080", nil)
}

func SearchCepUI(w http.ResponseWriter, r *http.Request) {
	cepParam := r.URL.Query().Get("cep")
	if cepParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cep, error := SearchCep(cepParam)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(os.Stdout).Encode(*cep)

	t := template.Must(template.New("template.html").ParseFiles("template.html"))

	err := t.Execute(w, cep)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing the request: %v\n", err)
	}
}

func SearchCepHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	cepParam := r.URL.Query().Get("cep")
	if cepParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cep, error := SearchCep(cepParam)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(*cep)
}

func SearchCep(cep string) (*ViaCEP, error) {
	var urlLeft string = "http://viacep.com.br/ws/"
	var urlRight string = "/json/"

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", urlLeft + cep + urlRight, nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var dataCep ViaCEP
	err = json.Unmarshal(body, &dataCep)

	return &dataCep, nil
}

func printDataCep(idxCep int, dataCep *ViaCEP) {
	fmt.Println(idxCep+1, "CEP:")
	fmt.Println("CEP: ", dataCep.Cep)
	fmt.Println("Logradouro: ", dataCep.Logradouro)
	fmt.Println("Bairro: ", dataCep.Bairro)
	fmt.Println("Localidade: ", dataCep.Localidade)
	fmt.Println("UF: ", dataCep.Uf)
	fmt.Println("IBGE: ", dataCep.Ibge)
	fmt.Println()
}

// We won't use this func anymore. Saving just in case I change idea to look for many ceps at the same time.
func SearchCeps() {
	var ceps []string
	numParams := len(os.Args) - 1

	// Populate the ceps array
	if numParams == 0 {
		var numCeps int
		var cep string

		println("How many CEPs you will check?")
		fmt.Scan(&numCeps)

		for i := 0; i < numCeps; i++ {
			fmt.Print("Enter the ", i+1, " CEP: ")
			fmt.Scan(&cep)
			ceps = append(ceps, cep)
		}

	} else {
		for _, cep := range os.Args[1:] {
			ceps = append(ceps, cep)
		}
	}

	for idxCep, cep := range ceps {
		dataCep, err := SearchCep(cep)
		if err != nil {

		}
		printDataCep(idxCep, dataCep)
	}
}
