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

func UserTerminalInput() models.UserInput {

	var userInput models.UserInput

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
	userInput.Token = userToken

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