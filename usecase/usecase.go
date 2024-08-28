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

func (u usecase) BulkData(token string){
	var row, succesTrans, failTrans int
	var totalPercobaan = 1
	uInput := helper.BulkDataTerminalInput()

	reseler := u.resellerRepo.GetResellerData()

	for i := 0; i < uInput.TotalInsertData; {

		user, err := u.custRepo.GetNIKFiltered(row, uInput.SheetChoose)
		if err != nil {
			return
		}
	
		if user.Code != uInput.TagSelected {
			row++
			continue
		}
		//change to db
		isAvailForTrans, isNewTransaction, currentTotalTrans, err := u.custRepo.GetHistoryTransactionDB(user.NIK, uInput.SheetChoose, uInput.UserMaxMonthPurchase, uInput.ItemPerPuchase)

		if err != nil {
			fmt.Println(err)
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
		trParam := helper.TransParamPrep(prData, userDetail, user.NIK, uInput.ItemPerPuchase)

		transResp, err := u.custRepo.CreateTransaction(trParam)

		if isNewTransaction {
			if err != nil {
				if errors.Is(err, helper.ErrTansFail){
					err = u.custRepo.CreateHistoryCustTransaction(user.NIK, 0, reseler.Data.DistrictName, user.Code, 0, uInput.SheetChoose)
					if err != nil {
						fmt.Println(err)
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
				fmt.Println(err)
				return	
			}
			err = u.custRepo.CreateHistoryCustTransaction(user.NIK, uInput.ItemPerPuchase, reseler.Data.DistrictName, user.Code, 1, uInput.SheetChoose)
			
			if err != nil {
				fmt.Println(err)	
				return
			}

		}

		if err != nil {
			if errors.Is(err, helper.ErrTansFail){
				err = u.custRepo.UpdateHistoryCustTranscation(user.NIK, uInput.SheetChoose, currentTotalTrans, 0)
				if err != nil{
					fmt.Print(err)
					return
				}
			}

				failTrans++
				totalPercobaan++
				fmt.Printf("Pesan: %s\n", transResp.Message)
				fmt.Println("Status: GAGAL")
				fmt.Println("========== TRANSAKSI SELESAI ==========")
				time.Sleep(25 * time.Second)
				continue	
			}
			
			err = u.custRepo.UpdateHistoryCustTranscation(user.NIK, uInput.SheetChoose, currentTotalTrans + uInput.ItemPerPuchase, 1)
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