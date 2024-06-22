package models

type ModeInput string

var ( 
	BULK ModeInput = "bulk"
	FILTER ModeInput = "filter"
)


type BulkInput struct {
	SheetChoose				string 	`json:"sheetChoose"`
	ColumnChoose			int		`json:"columnChoose"`
	TotalInsertData			int		`json:"totalInsertData"`
	TagSelected				string	`json:"tagSelected"`
	UserMaxMonthPurchase	int		`json:"userMaxMonthPurchase"`
}


type FilterDataInput struct {
	SheetChoose		string 	`json:"sheetChoose"`
	DataUpdate		int		`json:"dataUpdate"`
	ColumnChoose	int		`json:"columnChoose"`
}

type StartAppInput struct {
	Token 			string	`json:"token"`
	Mode 			string	`json:"mode"`
}