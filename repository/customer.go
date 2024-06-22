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
	"time"

	"github.com/xuri/excelize/v2"
)


type CustomerRepository interface {
	GetCustData(nationalId string) (models.GetCustomerResponse, error)
	CreateTransaction(param models.TransactionParam) (models.TransactionSuccesResponse, error)
	ReadCustsFromExcel(sheetChoose string, columnNumb int, index int) models.CustFromExcel
	WriteTransactionToExcel(userTrans models.CustToExcel, sheetChoose string) error
	GetRowsFiltered(kelurahan string) (int, error)
	ReadRowExcel(file string, sheet string, row int, col int) string
	WriteFilteredData(param models.WriteFilteredDataParam) error
	UpdateRowsFiltered(kelurahan string, numb int) error 
	GetNIKFiltered(row int, sheet string) (models.NIKFiltered, error)
	GetHistoryTransactionExcel(NIK string, sheet string, transMaks int) (bool, error) 
	UpdateCustHistoryTrans(sheet string, NIK string, keterangan string, tag string, isTransFail bool) (int, error)
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

func (u customerRepository) GetHistoryTransactionExcel(NIK string, sheet string, transMaks int) (bool, error) {
	f, err := excelize.OpenFile("libs/MAP_TRANSACTIONS.xlsx")
	
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer f.Close()

	sheetName := helper.FindOrCreateSheet(f, sheet)

	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Println(err)
		return false, err
	}

	for index, row := range rows {

		if len(rows) > 1 {

		if index > len(rows) {
			return true, err
			}
		
		if row[0] == NIK {
			if row[3] == "NO" {
				return false, err
				}
			num, err := strconv.Atoi(row[2])
			if err != nil {
				log.Println(err)
				return false, err
				}	
			if transMaks <= num {
				return false, err
				}
			}
		}
	}
	return true, nil
}

func (u customerRepository) GetCustData(nationalId string) (models.GetCustomerResponse, error) {
	var response models.GetCustomerResponse

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api-map.my-pertamina.id/customers/v1/verify-nik?nationalityId=%s", nationalId), nil)
	if err != nil {
		log.Println(err)
		return response, err
	}

	req.Header.Set("Authorization", u.token)
	req.Header.Set("Origin", "https://subsiditepatlpg.mypertamina.id")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	resp, err := u.httpClient.Do(req)

	if err != nil {
		log.Println("Errored when sending request to the server")
		return response, err
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return response, err
	}

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
        log.Println("Error decoding JSON:", err)
        return response, err
    }

	if response.Code == 429 {
		log.Println("To many request! Wait couple minutes")	
		return response, fmt.Errorf("ERROR")
	}

	if response.Code >= 400 && response.Code < 500 {
		if response.Code == 403 || response.Code == 401 {
			log.Println("Token Bermasalah!")
			return response, fmt.Errorf("ERROR")
		}
		return response, fmt.Errorf("ERROR")
	}
	return response, nil
}


func (u customerRepository) GetNIKFiltered(row int, sheet string) (models.NIKFiltered, error) {
	var nik models.NIKFiltered
	
	f, err := excelize.OpenFile("libs/DATA_FILTERED.xlsx")

	if err != nil {
		log.Println(err)
		return nik, err
	}
	defer f.Close()


	rows, err := f.GetRows(sheet)
	if err != nil {
		log.Fatalf("Sheet `%s` pada `%s` tidak ditemukan.", sheet, "libs/DATA_FILTERED.xlsx")
		return nik, err
	}

	nik.NIK = rows[row + 1][0]
	nik.Code = rows[row +1][1]
	return nik, nil	
}

func (u customerRepository) UpdateCustHistoryTrans(sheet string, NIK string, keterangan string, tag string, isTransFail bool) (int, error){
	var rowNumb int
	f, err := excelize.OpenFile("libs/MAP_TRANSACTIONS.xlsx")

	if err != nil{
		log.Println(err)
		return 0, err
	}	
	defer f.Close()

	sheetName := helper.FindOrCreateSheet(f, sheet)
	rows, err := f.GetRows(sheetName)
	if err != nil{
		log.Println(err)
		return 0, err
	}

	for index, row := range rows {
		if row[0] == NIK {
			if !isTransFail {
				f.SetCellValue(sheetName, "D"+strconv.Itoa(index+1), "NO")			
				f.SetCellValue(sheetName, "E"+strconv.Itoa(index+1), keterangan)
				f.Save()
				return -1, nil
			}
		    totalTrans, _ := strconv.Atoi(row[2])
			f.SetCellValue(sheetName, "C"+strconv.Itoa(index+1), strconv.Itoa(totalTrans + 1))
			f.SetCellValue(sheetName, "E"+strconv.Itoa(index+1), keterangan)
			f.Save()
			return 1, nil
		}

		if index + 1 == len(rows) {
			rowNumb = index + 1
			break
		} 
	}

	f.SetCellValue(sheetName, "A"+strconv.Itoa(rowNumb+1), NIK)
	f.SetCellValue(sheetName, "B"+strconv.Itoa(rowNumb+1), tag)
	f.SetCellValue(sheetName, "C"+strconv.Itoa(rowNumb+1), 1)
	f.SetCellValue(sheetName, "E"+strconv.Itoa(rowNumb+1), keterangan)
	if !isTransFail {
		f.SetCellValue(sheetName, "D"+strconv.Itoa(rowNumb+1), "NO")
	} else {
		f.SetCellValue(sheetName, "D"+strconv.Itoa(rowNumb+1), "YES")
	}
	f.Save()
	return 1, nil
}

func (u customerRepository) ReadCustsFromExcel(sheetChoose string, columnNumb int, index int) models.CustFromExcel{
	var custs models.CustFromExcel
	// var currentRow int

	// f, err := excelize.OpenFile("libs/MAP_TRANSACTIONS.xlsx")
	// defer f.Close()
	
	// if err != nil {
	// 	log.Println(err)
	// 	return custs
	// }

	// currentRow, _ = helper.FindSheetLength(f, sheetChoose)
	// if err := f.Save(); err != nil {
	// 	log.Println("Error saving MAP_TRANSACTIONS.xlsx:", err)
	// 	return custs
	// }	

	// if index == 0 {
	// 	helper.CheckNIK(custs.NIK)
	// }
	
	return custs
}


func (c customerRepository) ReadRowExcel(file string, sheet string, row int, col int) string{
	var data string
	
	f, err := excelize.OpenFile(file)
	if err != nil {
		log.Println(err)
		return data
	}

	rows, err := f.GetRows(sheet)
	if err != nil {
		log.Fatalf("Sheet `%s` pada `%s` tidak ditemukan.", sheet, file)
		return data	
	}

	data = rows[row][col]

	return data
}

func (c customerRepository) UpdateRowsFiltered(kelurahan string, numb int) error {
	f2, err := excelize.OpenFile("libs/DATA_FILTERED.xlsx")
	if err != nil {
		log.Println(err)
		return err
	}

	rows, err := f2.GetRows("home")
	if err != nil {
		log.Println(err)
		return err
	}

	for rowIndex, row := range rows {
		if len(row) >= 2 && row[0] == kelurahan {
			cell := fmt.Sprintf("B%d", rowIndex+1)
			f2.SetCellValue("home", cell, numb)

			if err := f2.Save(); err != nil {
				log.Println("Error saving DATA_FILTERED.xlsx:", err)
				return err
			}
		}
	}
	return nil
}

func (c customerRepository) GetRowsFiltered(kelurahan string) (int, error){

	f, err := excelize.OpenFile("libs/DATA_MAP_PANGKALAN_2024.xlsx")
	if err != nil {
		log.Println(err)
		return 0, err
	}

	_, err = f.GetRows(kelurahan)

	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer f.Close()

	f2, err := excelize.OpenFile("libs/DATA_FILTERED.xlsx")
	if err != nil {
		log.Println(err)
		return 0, err
	}

	rows, err := f2.GetRows("home")
	if err != nil {
		log.Println(err)
		return 0, err
	}

	for _, row := range rows {
		if len(row) >= 2 && row[0] == kelurahan {
			lastRowEdited, err := strconv.Atoi(row[1])
			if err != nil {
				return 0, err
			}
			return lastRowEdited, nil
		}
	}

	newRowIndex := len(rows) + 1
	f2.SetCellValue("home", fmt.Sprintf("A%d", newRowIndex), kelurahan)
	f2.SetCellValue("home", fmt.Sprintf("B%d", newRowIndex), 0)

	if err := f2.Save(); err != nil {
		log.Println("Error saving DATA_FILTERED.xlsx:", err)
		return 0, err
	}

	return 0, nil
}

func (u customerRepository) WriteFilteredData(param models.WriteFilteredDataParam) error{	
	f, err := excelize.OpenFile("libs/DATA_FILTERED.xlsx")

	if err != nil{
		log.Println(err)
		return err
	}

	sheetIndex, _ := f.GetSheetIndex(param.Sheet)
	if sheetIndex == -1 {
		f.NewSheet(param.Sheet)
		f.SetCellValue(param.Sheet, "A"+strconv.Itoa(1), "NIK")
		f.SetCellValue(param.Sheet, "B"+strconv.Itoa(1), "CODE")
		f.SetCellValue(param.Sheet, "C"+strconv.Itoa(1), "KETERANGAN")
		f.SetCellValue(param.Sheet, "D"+strconv.Itoa(1), "DATE FILTERED")
	}

	rows, err := f.GetRows(param.Sheet)
	if err != nil {
		return err
	}

	row := len(rows) + 1
	now := time.Now()
	formattedDate := now.Format("02/01/2006")
	code := helper.GetCustomerCode(param.Customer)

	f.SetCellValue(param.Sheet, "A"+strconv.Itoa(row), param.NIK)
	f.SetCellValue(param.Sheet, "B"+strconv.Itoa(row), code)
	f.SetCellValue(param.Sheet, "C"+strconv.Itoa(row), param.Keterangan)
	f.SetCellValue(param.Sheet, "D"+strconv.Itoa(row), formattedDate)

	if err := f.Save(); err != nil {
		log.Println("Error saving DATA_FILTERED.xlsx:", err)
		return err
	}

	return nil
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

	if errResponse.Code >= 400 && errResponse.Code < 500 {
		response.Code = errResponse.Code
		response.Message = errResponse.Message
		return response, helper.ErrTansFail
	}

	err = json.Unmarshal(responseBody, &response)
    if err != nil {
        log.Fatal("Error decoding JSON:", err)
        return response, err
    }

	return response, nil
}