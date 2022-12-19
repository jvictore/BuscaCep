# BuscaCep
This code was made to search information about CEPs (similar to postal code) in Brasil.

## There are 2 possibilities to use the application:


1 - You can pass the CEPs via command line arguments, example passing 2 CEPs:
  
    go run main.go 02402-020 91920-273
    // {02402-020 Rua Garção Tinoco  Santana São Paulo SP 3550308 1004 11 7107}
    // {91920-273 Rua Caliandra  Cavalhada Porto Alegre RS 4314902  51 8801}
    
2 - If you don't want to use the command line arguments, just run the application, and it will ask for the quantity of CEPs and the ones:

    go run main.go
    // How many CEPs you will check?
    2
    // Enter the 1 CEP: 
    02402-020
    // Enter the 2 CEP:
    91920-273
    // {02402-020 Rua Garção Tinoco  Santana São Paulo SP 3550308 1004 11 7107}
    // {91920-273 Rua Caliandra  Cavalhada Porto Alegre RS 4314902  51 8801}
