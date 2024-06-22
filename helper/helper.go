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
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return 0,""
	}
	startRow := len(rows) + 1

	return startRow, sheetName
}

func FindOrCreateSheet(f *excelize.File, sheetChoose string)(int, string){
	now := time.Now()
	_, week := now.ISOWeek()
	monthUpper := now.Format("Jan") 

	sheetName := fmt.Sprintf("%s-%s-WEEK%d", sheetChoose, monthUpper, week)

	numb, _ := f.GetSheetIndex(sheetName)
	if numb == -1 {
		f.NewSheet(sheetName)

		f.SetCellValue(sheetName, "A1", "Header 1")
		f.SetCellValue(sheetName, "B1", "Header 2")
		f.SetCellValue(sheetName, "C1", "Header 3")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return 0, ""
	}

	startRow := len(rows) + 1
	return startRow, sheetName
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


	userInput.SheetChoose = sheetChoose
	userInput.DataUpdate = totalUpdateData
	userInput.ColumnChoose = columnChoose

	return userInput
}

func BulkDataTerminalInput() models.BulkInput {

	var userInput models.BulkInput

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

	var userToken string
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

	clearTerminal()

	fmt.Println("Nilai sheet: ", sheetChoose)
	fmt.Println("Nilai index column: ", columnChoose)	
	fmt.Println("Nilai data diproses: ", totalInsertData)

	userInput.ColumnChoose = columnChoose
	userInput.SheetChoose = sheetChoose
	userInput.TotalInsertData =totalInsertData
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
		log.Fatal("Coba cek apakah index column sudah tepat")
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