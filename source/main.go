package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	// "net/smtp"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
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
}

func NewDataCep() *ViaCEP{
	return &ViaCEP{
		ID: uuid.New().String(),
		Cep: "",
		Logradouro: "",
		Bairro: "",
		Localidade: "",
		Uf: "",
		Ibge: "",
	}
}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/datacep")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	cep, err := SearchCep("13052200")
	if err != nil{
		panic(err)
	}
	erro := insertDataCep(db, cep)
	if erro != nil {
		panic(err)
	}
	// http.HandleFunc("/", SearchCepHandler)
	// http.HandleFunc("/ui", SearchCepUIHandler)
	// http.HandleFunc("/add", AddCepHandler)
	// http.ListenAndServe(":8080", nil)
}

func AddCepHandler(w http.ResponseWriter, r *http.Request){
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/datacep")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	cepParam := r.URL.Query().Get("cep")
	if cepParam == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	cep, err := SearchCep(cepParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	insertDataCep(db, cep)
}

func insertDataCep(db *sql.DB, cep *ViaCEP) error {
	stmt, err := db.Prepare("insert into dataceps(id, cep, logradouro, bairro, localidade, uf, ibge) values(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	
	_, err = stmt.Exec(cep.ID, cep.Cep, cep.Logradouro, cep.Bairro, cep.Localidade, cep.Uf, cep.Ibge)
	if err != nil {
		return err
	}

	return nil
}

func SearchCepUIHandler(w http.ResponseWriter, r *http.Request) {
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
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
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

	dataCep := NewDataCep()
	err = json.Unmarshal(body, dataCep)
	if err != nil {
		panic(err)
	}

	return dataCep, nil
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
