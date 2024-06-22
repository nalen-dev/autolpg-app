package main

import (
	"autolpg-app/helper"
	"autolpg-app/repository"
	"autolpg-app/usecase"
	"fmt"
	"net/http"
)

func main() {
	defer helper.ExitHandler()
	fmt.Println("Aplikasi telah berjalan.")

	client := &http.Client{}

	startUpInput := helper.StartAppTerminalInput()

	cr := repository.NewCustRepo(client, startUpInput.Token)
	pr := repository.NewProdRepo(client, startUpInput.Token)
	rr := repository.NewResellerRepo(client, startUpInput.Token)

	uc := usecase.CreateNewUseCase(cr, pr, rr)	


	switch startUpInput.Mode {
	case "bulk":
		uc.BulkData(startUpInput.Token)
		return;
	case "filtering":
		uc.FilteringData()
		return;
	}
}

