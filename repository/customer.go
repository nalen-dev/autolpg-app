package repository

import (
	"autolpg-app/helper"
	"autolpg-app/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/xuri/excelize/v2"
)


type CustomerRepository interface {
	GetCustData(nationalId string) (models.GetCustomerResponse,string)
	CreateTransaction(param models.TransactionParam) (models.TransactionSuccesResponse, error)
	ReadCustsFromExcel(sheetChoose string, columnNumb int, index int) models.CustFromExcel
	WriteTransactionToExcel(userTrans models.CustToExcel, sheetChoose string) error
}

type customerRepository struct {
	httpClient *http.Client
	token 		string
}

func NewCustRepo(httpClient *http.Client, token string) CustomerRepository{
	return &customerRepository{
			httpClient: httpClient, 
			token: token,
		}
}

func (u customerRepository) GetCustData(nationalId string) (models.GetCustomerResponse, string) {
	var response models.GetCustomerResponse
	var errResponse models.Response

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api-map.my-pertamina.id/customers/v1/verify-nik?nationalityId=%s", nationalId), nil)
	if err != nil {
		log.Fatal(err)
		return response, ""
	}

	req.Header.Set("Authorization", u.token)
	req.Header.Set("Origin", "https://subsiditepatlpg.mypertamina.id")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	resp, err := u.httpClient.Do(req)

	if err != nil {
		log.Println("Errored when sending request to the server")
		return response, ""
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

    err = json.Unmarshal(responseBody, &errResponse)
    if err != nil {
        log.Fatal("Error decoding JSON:", err)
        return response,""
    }

	if errResponse.Code == 429 {
		log.Fatal("To many request! Wait couple minutes")	
	}

	if errResponse.Code >= 400 && errResponse.Code < 500 {
		if errResponse.Code == 403 || errResponse.Code == 401 {
			log.Fatalf("Token Bermasalah!")
		}
		return response, "X"
	}


	err = json.Unmarshal(responseBody, &response)
	if err != nil {
        log.Fatal("Error decoding JSON:", err)
        return response, ""
    }
	return response, ""
}

func (u customerRepository) ReadCustsFromExcel(sheetChoose string, columnNumb int, index int) models.CustFromExcel{
	var custs models.CustFromExcel
	var currentRow int

	f, err := excelize.OpenFile("libs/MAP_TRANSACTIONS.xlsx")
	if err != nil {
		log.Println(err)
		return custs
	}

	currentRow, _ = helper.FindSheetLength(f, sheetChoose)
	f.Close()

	f, err = excelize.OpenFile("libs/DATA_MAP_PANGKALAN_2024.xlsx")
	if err != nil {
		log.Println(err)
		return custs
	}

	rows, err := f.GetRows(sheetChoose)
	if err != nil {
		log.Fatalf("Sheet `%s` pada DATA_MAP_PANGKALAN_2024.xlsx tidak ditemukan.", sheetChoose)
		return custs	
	}

	custs.NumbRow = currentRow - 1
	custs.NIK = rows[currentRow - 1][columnNumb]
	
	if index == 0 {
		helper.CheckNIK(custs.NIK)
	}
	
	return custs
}

func (u customerRepository) WriteTransactionToExcel(cust models.CustToExcel, sheetChoose string) error {
	filePath := "libs/MAP_TRANSACTIONS.xlsx"
	var f *excelize.File
	var err error

	if f, err = excelize.OpenFile(filePath); err != nil {
		f = excelize.NewFile()
	}

	currentRow, sheetName := helper.FindSheetLength(f, sheetChoose)

	f.SetCellValue(sheetName, "A"+strconv.Itoa(currentRow), cust.NumbRow)
	f.SetCellValue(sheetName, "B"+strconv.Itoa(currentRow), cust.NIK)
	f.SetCellValue(sheetName, "C"+strconv.Itoa(currentRow), cust.CAT)
	f.SetCellValue(sheetName, "D"+strconv.Itoa(currentRow), cust.TransactionId)
	f.SetCellValue(sheetName, "E"+strconv.Itoa(currentRow), cust.Status)
	


	// Simpan perubahan ke file Excel
	if err := f.SaveAs(filePath); err != nil {
		return err
	}

	return nil
}

func (u customerRepository) CreateTransaction(param models.TransactionParam) (models.TransactionSuccesResponse, error) {

	var response models.TransactionSuccesResponse
	var errResponse models.Response

	jsonData, err := json.Marshal(param)
	if err != nil {
		log.Println(err)
		return response, err
	}

	req, err := http.NewRequest(http.MethodPost,"https://api-map.my-pertamina.id/general/v1/transactions", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		return response, err
	}

	req.Header.Set("Authorization", u.token)
	req.Header.Set("Origin", "https://subsiditepatlpg.mypertamina.id")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	req.Header.Set("Content-Type","application/json")

	resp, err := u.httpClient.Do(req)
	
	if err != nil {
		log.Fatalf("Errored when sending request to the server: %s", err)
		return response, err
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return response, err
	}

	err = json.Unmarshal(responseBody, &errResponse)
    if err != nil {
        log.Fatal("Error decoding JSON:", err)
        return response, err
    }

	if errResponse.Code >=400 && errResponse.Code < 500 {
		response.Code = errResponse.Code
		response.Message = errResponse.Message
		return response, nil 
	}

	err = json.Unmarshal(responseBody, &response)
    if err != nil {
        log.Fatal("Error decoding JSON:", err)
        return response, err
    }

	return response, err
}