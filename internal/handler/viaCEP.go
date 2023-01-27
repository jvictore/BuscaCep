package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	// _ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"

	models "github.com/jvictore/ZipCodeFinder/internal/models"
)

func UpdateCepHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/datacep")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	cepParam := r.URL.Query().Get("cep")
	if cepParam != "" {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func AddCepHandler(w http.ResponseWriter, r *http.Request) {
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

	InsertDataCep(db, cep)
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

func InsertDataCep(db *sql.DB, cep *models.ViaCEP) error {
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

func SearchCep(cep string) (*models.ViaCEP, error) {
	var urlLeft string = "http://viacep.com.br/ws/"
	var urlRight string = "/json/"

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", urlLeft+cep+urlRight, nil)
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

func NewDataCep() *models.ViaCEP {
	return &models.ViaCEP{
		ID:         uuid.New().String(),
		Cep:        "",
		Logradouro: "",
		Bairro:     "",
		Localidade: "",
		Uf:         "",
		Ibge:       "",
	}
}
