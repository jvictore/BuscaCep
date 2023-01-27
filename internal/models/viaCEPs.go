package models

type ViaCEP struct {
	ID         string
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	Uf         string `json:"uf"`
	Ibge       string `json:"ibge"`
}
