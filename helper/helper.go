package helper

import (
	"autolpg-app/models"
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

func FindSheetLength(f *excelize.File, sheetChoose string) (int, string){

	now := time.Now()
	_, week := now.ISOWeek()
	monthUpper := now.Format("Jan") // Mendapatkan tiga huruf pertama dari bulan

	sheetName := fmt.Sprintf("%s-%s-WEEK%d", sheetChoose, monthUpper, week)

	numb, _ := f.GetSheetIndex(sheetName)
	if numb == -1 {
		f.NewSheet(sheetName)
		f.SetCellValue(sheetName, "A1", "NIK")
		f.SetCellValue(sheetName, "B1", "TAG")
		f.SetCellValue(sheetName, "C1", "TRANSACTION ID")
		f.SetCellValue(sheetName, "D1", "Jumlah Transaksi Bulanan")
		f.SetCellValue(sheetName, "E1", "KETERANGAN")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return 0,""
	}
	startRow := len(rows) + 1

	return startRow, sheetName
}

func FindOrCreateSheet(f *excelize.File, sheetChoose string)(string, error){
	t := time.Now()
	firstOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())

	weekday := int(firstOfMonth.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	dayOfMonth := t.Day()
	weekOfMonth := (dayOfMonth-1+weekday-1)/7 + 1
	monthUpper := firstOfMonth.Format("Jan")
	sheetName := fmt.Sprintf("%s-%s-WEEK%d", sheetChoose, monthUpper, weekOfMonth)

	numb, _ := f.GetSheetIndex(sheetName)
	if numb == -1 {
		f.NewSheet(sheetName)
		f.SetCellValue(sheetName, "A1", "NIK")
		if err := f.MergeCell(sheetName, "A1", "A2"); err != nil {
			log.Printf("Error merging cells: %v\n", err)
			return "", err
		}

		f.SetCellValue(sheetName, "B1", "TAG")
		if err := f.MergeCell(sheetName, "B1", "B2"); err != nil {
			log.Printf("Error merging cells: %v\n", err)
			return "", err
		}
		
		f.SetCellValue(sheetName, "C1", "Jumlah Transaksi Bulanan")
		if err := f.MergeCell(sheetName, "C1", "C2"); err != nil {
			log.Printf("Error merging cells: %v\n", err)
			return "", err
		}

		f.SetCellValue(sheetName, "D1", "isAvailable")
		if err := f.MergeCell(sheetName, "D1", "D2"); err != nil {
			log.Printf("Error merging cells: %v\n", err)
			return "", err
		}

		f.SetCellValue(sheetName, "E1", "KETERANGAN")
		if err := f.MergeCell(sheetName, "E1", "E2"); err != nil {
			log.Printf("Error merging cells: %v\n", err)
			return "", err
		}

		f.SetCellValue(sheetName, "G1", "TOTAL")
		if err := f.MergeCell(sheetName, "G1", "G2"); err != nil {
			log.Printf("Error merging cells: %v\n", err)
			return "", err
		}
		
		f.SetCellValue(sheetName, "H1", "RT")
		f.SetCellFormula(sheetName, "H2", `SUMIF(B3:B9999, "RT", C3:C9999)`)

		f.SetCellValue(sheetName, "I1", "UM")
		f.SetCellFormula(sheetName, "I2", `SUMIF(B3:B9999, "UM", C3:C9999)`)


		style, err := f.NewStyle(&excelize.Style{Alignment: &excelize.Alignment{Horizontal:"center", Vertical: "center"}})
		if err != nil {
			log.Fatalf("Error creating style: %v", err)
		}

		f.SetCellStyle(sheetName, "A1", "A2", style)
		f.SetCellStyle(sheetName, "B1", "B2", style)
		f.SetCellStyle(sheetName, "C1", "C2", style)
		f.SetCellStyle(sheetName, "D1", "D2", style)
		f.SetCellStyle(sheetName, "E1", "E2", style)
		f.SetCellStyle(sheetName, "F1", "F2", style)
		f.SetCellStyle(sheetName, "G1", "G2", style)

		f.Save()
	}
	return sheetName, nil
}

func GetCustomerCode(cust models.Customer) string {
		switch cust.CustomerTypes[len(cust.CustomerTypes) - 1].Name {
		case "Rumah Tangga":
			return "RT"
		case "Usaha Mikro":
			return "UM"
		default:
			return "PE"			
		}
}

func StartAppTerminalInput() models.StartAppInput {

	var userToken string
	var startInput models.StartAppInput

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Masukkan token: ")
		if scanner.Scan() {
			userToken = scanner.Text()
		}
		if userToken != "" {
			break
		} else {
			fmt.Println("Nilai tidak boleh kosong. Silakan masukkan kembali.")
		}
	}

	var mode string
	for {
        fmt.Println("Pilih mode yang anda inginkan:")
        fmt.Println("1. Bulk Insert")
        fmt.Println("2. Filtering Data")
        fmt.Print("Masukkan nomor mode (1 atau 2): ")

        if scanner.Scan() {
            mode = scanner.Text()
        }

        if mode == "1" || mode == "2" {
			
			if mode == "1" {
				startInput.Mode = "bulk"
			}

			if mode == "2" {
				startInput.Mode = "filtering"
			}

            break
        } else {
            fmt.Println("Input tidak valid. Silakan masukkan nomor 1 atau 2.")
        }
    }
	clearTerminal()

	startInput.Token = userToken
	return startInput
}


func FilterDataTerminalInpit() models.FilterDataInput {
	var userInput models.FilterDataInput
	scanner := bufio.NewScanner(os.Stdin)

	var sheetChoose string
	for {
		fmt.Print("Masukkan sheet name yang diinginkan pada file DATA_MAP_PANGKALAN_2024.xlsx: ")
		if scanner.Scan() {
			sheetChoose = scanner.Text()
		}

		if sheetChoose != "" {
			break
		} else {
			fmt.Println("Nilai tidak boleh kosong. Silakan masukkan kembali.")
		}
	}

	var totalUpdateData int
	for {
		fmt.Print("Masukkan jumlah data yang ingin diproses: ")
		if scanner.Scan() {
			input := scanner.Text()
			num, err := strconv.Atoi(input)
			if err != nil {
				log.Printf("Input tidak valid: %v\n", err)
				continue
			}
			if num > 0 {
				totalUpdateData = num
				break
			} else {
				fmt.Println("Total update data harus lebih dari 0. Silakan masukkan kembali.")
			}
		}
	}

	var columnChoose int
	for {
		fmt.Printf("Masukkan index column NIK pada sheet %s: ", sheetChoose)
		if scanner.Scan() {
			input := scanner.Text()
		
			num, err := strconv.Atoi(input)
			if err != nil {
				log.Printf("Input tidak valid untuk totalInsertData: %v\n", err)
				continue
			}
			if num >= 0 {
				columnChoose = num
				break
			}
		}
	}

	clearTerminal()

	userInput.SheetChoose = sheetChoose
	userInput.DataUpdate = totalUpdateData
	userInput.ColumnChoose = columnChoose
	fmt.Println("Nilai sheet: ", sheetChoose)
	fmt.Println("Nilai index column: ", columnChoose)	
	fmt.Println("Nilai data diproses: ", totalUpdateData)

	return userInput
}

func BulkDataTerminalInput() models.BulkInput {

	var userInput models.BulkInput

	scanner := bufio.NewScanner(os.Stdin)

	var sheetChoose string
	for {
		fmt.Print("Masukkan sheet name yang diinginkan pada file DATA_FILTERED.xlsx: ")
		if scanner.Scan() {
			sheetChoose = scanner.Text()
		}

		if sheetChoose != "" {
			break
		} else {
			fmt.Println("Nilai tidak boleh kosong. Silakan masukkan kembali.")
		}
	}

	var tagSelected string
	for {
        fmt.Println("Pilih kategori customer yang anda inginkan:")
        fmt.Println("1. RT (Rumah Tangga)")
        fmt.Println("2. UM (Usaha Mikro)")
        fmt.Print("Masukkan nomor category (1 atau 2): ")

        if scanner.Scan() {
            tagSelected = scanner.Text()
        }

        if tagSelected == "1" || tagSelected == "2" {
			
			if tagSelected == "1" {
				userInput.TagSelected = "RT"
			}

			if tagSelected == "2" {
				userInput.TagSelected = "UM"
			}

            break
        } else {
            fmt.Println("Input tidak valid. Silakan masukkan nomor 1 atau 2.")
        }
    }

	var maxPurchase int
	for {
		fmt.Print("Berapa banyak customer bisa beli dalam 1 bulan?: ")
		if scanner.Scan() {
			input := scanner.Text()
			num, err := strconv.Atoi(input)
			if err != nil {
				log.Printf("Input tidak valid untuk max purchase: %v\n", err)
				continue
			}
			if num > 0 {
				maxPurchase = num
				break
			} else {
				fmt.Println("maxPurchase harus lebih dari 0. Silakan masukkan kembali.")
			}
		}
	}

	var totalInsertData int
	for {
		fmt.Print("Masukkan jumlah data yang ingin diproses: ")
		if scanner.Scan() {
			input := scanner.Text()
			num, err := strconv.Atoi(input)
			if err != nil {
				log.Printf("Input tidak valid untuk totalInsertData: %v\n", err)
				continue
			}
			if num > 0 {
				totalInsertData = num
				break
			} else {
				fmt.Println("TotalInsertData harus lebih dari 0. Silakan masukkan kembali.")
			}
		}
	}

	clearTerminal()

	fmt.Println("Nilai sheet: ", sheetChoose)
	fmt.Println("Nilai customer category dipilih: ", userInput.TagSelected)	
	fmt.Println("Max Purchase: ", maxPurchase)
	fmt.Println("Nilai data diproses: ", totalInsertData)

	userInput.ColumnChoose = 0
	userInput.SheetChoose = sheetChoose
	userInput.TotalInsertData =totalInsertData
	userInput.UserMaxMonthPurchase = maxPurchase
	return userInput
}

func clearTerminal() {
	var cmd *exec.Cmd

	if osname := os.Getenv("GOOS"); osname == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}

func CheckNIK(nik string) {
	regex := regexp.MustCompile(`^\d{16}$`)
	if regex.MatchString(nik) {
	} else {
		log.Println("Coba cek apakah index column sudah tepat")
	}
}

func ExitHandler() {
    fmt.Println("Press Enter to exit...")
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        fmt.Println("Exiting...")
        break
    }
}

func TransParamPrep(prData models.GetProdResponse, userDetail models.GetCustomerResponse, NIK string) models.TransactionParam {
		var transactionParam models.TransactionParam

		transactionParam.Products = make([]models.Products, 1)
		transactionParam.Products[0].ProductID = prData.Data.ProductID
		transactionParam.Products[0].Quantity = 1
		transactionParam.InputNominal = prData.Data.Price
		transactionParam.Change = 0
		transactionParam.PaymentType = "cash"
		transactionParam.Subsidi.NIK = NIK
		transactionParam.Subsidi.FamilyID = userDetail.Data.FamilyId
		transactionParam.Subsidi.Category = userDetail.Data.CustomerTypes[len(userDetail.Data.CustomerTypes)-1].Name
		transactionParam.Subsidi.SourceTypeID = userDetail.Data.CustomerTypes[len(userDetail.Data.CustomerTypes)-1].SourceTypeId
		transactionParam.Subsidi.Nama = userDetail.Data.Name
		transactionParam.Subsidi.ChannelInject = userDetail.Data.ChannelInject
		transactionParam.Subsidi.PengambilanItemSubsidi = make([]models.PengambilanItemSubsidi, 1)
		transactionParam.Subsidi.PengambilanItemSubsidi[0].Item = "ELPIJI"
		transactionParam.Subsidi.PengambilanItemSubsidi[0].PotonganHarga = 0
		transactionParam.Subsidi.PengambilanItemSubsidi[0].Quantitas = 1
		
		return transactionParam
}