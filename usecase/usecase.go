package usecase

import (
	"autolpg-app/helper"
	"autolpg-app/models"
	"autolpg-app/repository"
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


func (u usecase) BulkData(token string){
	uInput := helper.BulkDataTerminalInput()

	for i := 0; i < uInput.TotalInsertData; i++ {
		var userCode string
		
		data := u.custRepo.ReadCustsFromExcel(uInput.SheetChoose, uInput.ColumnChoose, i)
		log.Printf("Memproses data NIK: %s\n", data.NIK)
		
		userData, userCode := u.custRepo.GetCustData(data.NIK)
		if userCode != "X" {
			prData := u.prodRepo.GetProductData()
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
		 transaction.Subsidi.Category = userData.Data.CustomerTypes[len(userData.Data.CustomerTypes)-1].Name
		 transaction.Subsidi.SourceTypeID = userData.Data.CustomerTypes[len(userData.Data.CustomerTypes)-1].SourceTypeId
		 transaction.Subsidi.Nama = userData.Data.Name
		 transaction.Subsidi.ChannelInject = userData.Data.ChannelInject
		 transaction.Subsidi.PengambilanItemSubsidi = make([]models.PengambilanItemSubsidi, 1)
		 transaction.Subsidi.PengambilanItemSubsidi[0].Item = "ELPIJI"
		 transaction.Subsidi.PengambilanItemSubsidi[0].PotonganHarga = 0
		 transaction.Subsidi.PengambilanItemSubsidi[0].Quantitas = 1
		 resp, _ := u.custRepo.CreateTransaction(transaction)
		 
		if resp.Code >= 400 && resp.Code < 500 {
			u.custRepo.WriteTransactionToExcel(models.CustToExcel{ 
				NumbRow: data.NumbRow, 
				NIK: data.NIK, 
				CAT: userCode, 
				TransactionId: "-",
				Status: resp.Message,
				}, uInput.SheetChoose)	
		}

		if resp.Code >= 200 && resp.Code < 300 {
			u.custRepo.WriteTransactionToExcel(models.CustToExcel{ 
				NumbRow: data.NumbRow, 
				NIK: data.NIK, 
				CAT: userCode, 
				TransactionId: resp.Data.TransactionId,
				Status: "Sukses",
				}, uInput.SheetChoose)			
		}
		} else {
			u.custRepo.WriteTransactionToExcel(models.CustToExcel{
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


func(u usecase) FilteringData(){

	userInput := helper.FilterDataTerminalInpit()

	lastRow, _ := u.custRepo.GetRowsFiltered(userInput.SheetChoose)

	for i := lastRow + 1; i < userInput.DataUpdate + lastRow + 1; i++ {
		var insertFilteredData models.WriteFilteredDataParam
	
		NIK := u.custRepo.ReadRowExcel("libs/DATA_MAP_PANGKALAN_2024.xlsx", userInput.SheetChoose, i, userInput.ColumnChoose)
		
		log.Printf("NIK %s diproses", NIK)
		s1, e1 := u.custRepo.GetCustData(NIK)
		
		if e1 != "X" {
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