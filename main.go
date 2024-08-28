package main

import (
	"autolpg-app/helper"
	"autolpg-app/repository"
	"autolpg-app/usecase"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	defer helper.ExitHandler()
	fmt.Println("Aplikasi telah berjalan.")

	db, err := sql.Open("sqlite3", "./libs/autolpg.db");
	if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
	fmt.Println("Berhasil connect ke database")

	client := &http.Client{}

	startUpInput := helper.StartAppTerminalInput()

	cr := repository.NewCustRepo(client, startUpInput.Token, db)
	pr := repository.NewProdRepo(client, startUpInput.Token, db)
	rr := repository.NewResellerRepo(client, startUpInput.Token, db)

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

