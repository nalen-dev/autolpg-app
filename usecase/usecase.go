package usecase

import (
	"autolpg-app/helper"
	"autolpg-app/models"
	"autolpg-app/repository"
	"errors"
	"fmt"
	"log"
	"time"
)

type Usecase interface{
	BulkData(token string)
	FilteringData()
}

type usecase struct {
	custRepo 		repository.CustomerRepository
	prodRepo 	 	repository.ProductRepository
	resellerRepo 	repository.ResellerRepository
}

func CreateNewUseCase(custRepo repository.CustomerRepository, prodRepo repository.ProductRepository, resellerRepo repository.ResellerRepository	) Usecase {
	return &usecase{
			custRepo: custRepo,
			prodRepo: prodRepo,
			resellerRepo: resellerRepo,
	}
}

// func (u usecase) BulkData(token string){
// 	uInput := helper.BulkDataTerminalInput()

// 	for i := 0; i < uInput.TotalInsertData; i++ {
// 		var userCode string
		
// 		data := u.custRepo.ReadCustsFromExcel(uInput.SheetChoose, uInput.ColumnChoose, i)
// 		log.Printf("Memproses data NIK: %s\n", data.NIK)
		
// 		userData, userCode := u.custRepo.GetCustData(data.NIK)
// 		if userCode != "X" {
// 			prData := u.prodRepo.GetProductData()
// 			userCode =	helper.GetCustomerCode(userData.Data)
// 		 var transaction models.TransactionParam

// 		 transaction.Products = make([]models.Products, 1)
// 		 transaction.Products[0].ProductID = prData.Data.ProductID
// 		 transaction.Products[0].Quantity = 1
// 		 transaction.InputNominal = prData.Data.Price
// 		 transaction.Change = 0
// 		 transaction.PaymentType = "cash"
// 		 transaction.Subsidi.NIK = data.NIK
// 		 transaction.Subsidi.FamilyID = userData.Data.FamilyId
// 		 transaction.Subsidi.Category = userData.Data.CustomerTypes[len(userData.Data.CustomerTypes)-1].Name
// 		 transaction.Subsidi.SourceTypeID = userData.Data.CustomerTypes[len(userData.Data.CustomerTypes)-1].SourceTypeId
// 		 transaction.Subsidi.Nama = userData.Data.Name
// 		 transaction.Subsidi.ChannelInject = userData.Data.ChannelInject
// 		 transaction.Subsidi.PengambilanItemSubsidi = make([]models.PengambilanItemSubsidi, 1)
// 		 transaction.Subsidi.PengambilanItemSubsidi[0].Item = "ELPIJI"
// 		 transaction.Subsidi.PengambilanItemSubsidi[0].PotonganHarga = 0
// 		 transaction.Subsidi.PengambilanItemSubsidi[0].Quantitas = 1
// 		 resp, _ := u.custRepo.CreateTransaction(transaction)
		 
// 		if resp.Code >= 400 && resp.Code < 500 {
// 			u.custRepo.WriteTransactionToExcel(models.CustToExcel{ 
// 				NumbRow: data.NumbRow, 
// 				NIK: data.NIK, 
// 				CAT: userCode, 
// 				TransactionId: "-",
// 				Status: resp.Message,
// 				}, uInput.SheetChoose)	
// 		}

// 		if resp.Code >= 200 && resp.Code < 300 {
// 			u.custRepo.WriteTransactionToExcel(models.CustToExcel{ 
// 				NumbRow: data.NumbRow, 
// 				NIK: data.NIK, 
// 				CAT: userCode, 
// 				TransactionId: resp.Data.TransactionId,
// 				Status: "Sukses",
// 				}, uInput.SheetChoose)			
// 		}
// 		} else {
// 			u.custRepo.WriteTransactionToExcel(models.CustToExcel{
// 				NumbRow: data.NumbRow, 
// 				NIK: data.NIK, 
// 				CAT: userCode,
// 				TransactionId: "-",
// 				Status: "NIK tidak valid/tidak terdaftar",
// 				}, uInput.SheetChoose)
		
// 		}
// 		log.Println("----------- Data berhasil diproses ----------------")
// 		time.Sleep(30 * time.Second)
// 	}

// 	log.Printf("Selesai, %d NIK berhasil diproses", uInput.TotalInsertData)
// }


func (u usecase) BulkData(token string){
	var row, succesTrans, failTrans int
	var totalPercobaan = 1
	uInput := helper.BulkDataTerminalInput()

	for i := 0; i < uInput.TotalInsertData; {
		user, err := u.custRepo.GetNIKFiltered(row, uInput.SheetChoose)
		if err != nil {
			return
		}
	
		if user.Code != uInput.TagSelected {
			row++
			continue
		}
	
		isAvailForTrans, err := u.custRepo.GetHistoryTransactionExcel(user.NIK, uInput.SheetChoose, uInput.UserMaxMonthPurchase)

		if err != nil {
			return 
		}

		if !isAvailForTrans {
			row++
			continue
		}

		userDetail, err := u.custRepo.GetCustData(user.NIK)

		if err != nil {
			fmt.Println(userDetail.Message)
			return 
		}

		log.Printf("\n========== PERCOBAAN TRANSAKSI ke-%d ==========\nNIK: %s\n", totalPercobaan, user.NIK)
		prData := u.prodRepo.GetProductData()
		trParam := helper.TransParamPrep(prData, userDetail, user.NIK)

		transResp, err := u.custRepo.CreateTransaction(trParam)
		if err != nil {
			if errors.Is(err, helper.ErrTansFail){
				 _, err := u.custRepo.UpdateCustHistoryTrans(uInput.SheetChoose, user.NIK, transResp.Message, user.Code, false)
				if err != nil {
					return
				}
				failTrans++
				totalPercobaan++
				
				fmt.Printf("Pesan: %s\n", transResp.Message)
				fmt.Println("Status: GAGAL")
				fmt.Println("========== TRANSAKSI SELESAI ==========")
				time.Sleep(25 * time.Second)
				continue
			}
			return 
		}
		_, err = u.custRepo.UpdateCustHistoryTrans(uInput.SheetChoose, user.NIK, transResp.Message, user.Code, true)
		if err != nil {
			return
		}
		succesTrans++
		i++
		row++
		totalPercobaan++
		fmt.Printf("Pesan: %s\n", transResp.Message)
		fmt.Println("Status: BERHASIL")
		fmt.Println("========== TRANSAKSI SELESAI ==========")
		time.Sleep(25 * time.Second)
	}

	log.Printf("\nSelesai :\n%d Transaksi berhasil diproses\n%d Transaksi gagal\n", succesTrans, failTrans)
}


func(u usecase) FilteringData(){

	userInput := helper.FilterDataTerminalInpit()

	lastRow, _ := u.custRepo.GetRowsFiltered(userInput.SheetChoose)

	for i := lastRow + 1; i < userInput.DataUpdate + lastRow + 1; i++ {
		var insertFilteredData models.WriteFilteredDataParam
	
		NIK, err := u.custRepo.ReadRowExcel("libs/DATA_MAP_PANGKALAN_2024.xlsx", userInput.SheetChoose, i, userInput.ColumnChoose)
		
		if err != nil{
			return
		}

		log.Printf("NIK %s diproses", NIK)
		s1, err := u.custRepo.GetCustData(NIK)
		
		if err == nil {
			insertFilteredData.Sheet = userInput.SheetChoose
			insertFilteredData.NIK = NIK
			insertFilteredData.Customer = s1.Data
			insertFilteredData.Keterangan = "Sukses"
			if err := u.custRepo.WriteFilteredData(insertFilteredData); err != nil{
				return
			}
		} 

		if s1.Code == 429 {
			return
		}

		 if err := u.custRepo.UpdateRowsFiltered(userInput.SheetChoose, i+1); err != nil {
			return
		 }

		log.Println("------- Berhasil Memproses Data --------")
		time.Sleep(30 * time.Second)
	}

}