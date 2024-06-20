package models

type UserInput struct {
	Token 			string	`json:"token"`
	SheetChoose		string 	`json:"sheetChoose"`
	ColumnChoose	int		`json:"columnChoose"`
	TotalInsertData	int		`json:"totalInsertData"`
}