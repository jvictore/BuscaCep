# Zip Code Finder
This code was made to search information about CEPs (similar to zip code) in Brasil.

## There are 2 possibilities to use the application:


1 - You can pass the CEPs via command line arguments, example passing 2 CEPs:
  
    go run main.go 02402-020 91920-273
    
    // 1 CEP:
    // CEP:  02402-020
    // Logradouro:  Rua Garção Tinoco
    // Complemento:  
    // Bairro:  Santana
    // Localidade:  São Paulo
    // UF:  SP
    // IBGE:  3550308
    
    // 2 CEP:
    // CEP:  91920-273
    // Logradouro:  Rua Caliandra
    // Complemento:  
    // Bairro:  Cavalhada
    // Localidade:  Porto Alegre
    // UF:  RS
    // IBGE:  4314902
    
2 - If you don't want to use the command line arguments, just run the application, and it will ask for the quantity of CEPs and the ones:

    go run main.go
    
    // How many CEPs you will check?
    2
    
    // Enter the 1 CEP: 
    02402-020
    
    // Enter the 2 CEP:
    91920-273
        
    // 1 CEP:
    // CEP:  02402-020
    // Logradouro:  Rua Garção Tinoco
    // Complemento:  
    // Bairro:  Santana
    // Localidade:  São Paulo
    // UF:  SP
    // IBGE:  3550308
    
    // 2 CEP:
    // CEP:  91920-273
    // Logradouro:  Rua Caliandra
    // Complemento:  
    // Bairro:  Cavalhada
    // Localidade:  Porto Alegre
    // UF:  RS
    // IBGE:  4314902

CREATE TABLE dataceps (id varchar(255), cep varchar(9) UNIQUE, logradouro varchar(255), bairro varchar(255), localidade varchar(255), uf varchar(2), ibge varchar(255), primary key(id));
