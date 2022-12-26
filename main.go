package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"
)

type ViaCEP struct {
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

	c := http.Client{Timeout: time.Second}
	req, err := c.Get(urlLeft + cep + urlRight)
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
