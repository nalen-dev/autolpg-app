package main

import (
	"autolpg-app/helper"
	"autolpg-app/models"
	"autolpg-app/repository"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Aplikasi telah berjalan.")

	uInput := helper.UserTerminalInput()

	client := &http.Client{}

	cr := repository.NewCustRepo(client, uInput.Token)
	pr := repository.NewProdRepo(client, uInput.Token)
	for i := 0; i < uInput.TotalInsertData; i++ {
		var userCode string
		
		data := cr.ReadCustsFromExcel(uInput.SheetChoose, uInput.ColumnChoose, i)
		log.Printf("Memproses data NIK: %s\n", data.NIK)
		
		userData, userCode := cr.GetCustData(data.NIK)
		if userCode != "X" {
			prData := pr.GetProductData()
			userCode =	helper.GetCustomerCode(userData.Data)
		 var transaction models.TransactionParam

		 transaction.Products = make([]models.Products, 1)
		 transaction.Products[0].ProductID = prData.Data.ProductID
		 transaction.Products[0].Quantity = 1
		 transaction.InputNominal = prData.Data.Price
		 transaction.Change = 0
		 transaction.PaymentType = "cash"
		 transaction.Subsidi.NIK = data.NIK
		 transaction.Subsidi.FamilyID = userData.Data.FamilyId
		 transaction.Subsidi.Category = userData.Data.CustomerTypes[0].Name
		 transaction.Subsidi.SourceTypeID = userData.Data.CustomerTypes[0].SourceTypeId
		 transaction.Subsidi.Nama = userData.Data.Name
		 transaction.Subsidi.ChannelInject = userData.Data.ChannelInject
		 transaction.Subsidi.PengambilanItemSubsidi = make([]models.PengambilanItemSubsidi, 1)
		 transaction.Subsidi.PengambilanItemSubsidi[0].Item = "ELPIJI"
		 transaction.Subsidi.PengambilanItemSubsidi[0].PotonganHarga = 0
		 transaction.Subsidi.PengambilanItemSubsidi[0].Quantitas = 1
		 resp, _ := cr.CreateTransaction(transaction)
		 
		if resp.Code >= 400 && resp.Code < 500 {
			cr.WriteTransactionToExcel(models.CustToExcel{ 
				NumbRow: data.NumbRow, 
				NIK: data.NIK, 
				CAT: userCode, 
				TransactionId: "-",
				Status: resp.Message,
				}, uInput.SheetChoose)	
		}

		if resp.Code >= 200 && resp.Code < 300 {
			cr.WriteTransactionToExcel(models.CustToExcel{ 
				NumbRow: data.NumbRow, 
				NIK: data.NIK, 
				CAT: userCode, 
				TransactionId: resp.Data.TransactionId,
				Status: "Sukses",
				}, uInput.SheetChoose)			
		}
		} else {
			cr.WriteTransactionToExcel(models.CustToExcel{
				NumbRow: data.NumbRow, 
				NIK: data.NIK, 
				CAT: userCode,
				TransactionId: "-",
				Status: "NIK tidak valid/tidak terdaftar",
				}, uInput.SheetChoose)
		
		}
		log.Println("----------- Data berhasil diproses ----------------")
		time.Sleep(30 * time.Second)
	}

	log.Printf("Selesai, %d NIK berhasil diproses", uInput.TotalInsertData)

}

