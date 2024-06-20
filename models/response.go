package models

type Response struct {
    Success bool `json:"success"`
    Code    int  `json:"code"`
    Message string `json:"message"`
}

type GetCustomerResponse struct {
	Response
	Data Customer `json:"data"`
}

type GetProdResponse struct {
	Response
	Data Product `json:"data"`
}

type GetResellerResponse struct {
	Response
	Data Reseller `json:"data"`
}

type TransactionSuccesResponse struct {
	Response
	Data TransactionSucces 	`json:"data"`
}