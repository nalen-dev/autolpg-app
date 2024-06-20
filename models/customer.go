package models

type CustFromExcel struct {
	NumbRow 	int `json:"row"`
	NIK		 	string  `json:"NIK"`
}

type CustToExcel struct {
	NumbRow 		int 	`json:"row"`
	NIK		 		string  `json:"NIK"`
	CAT				string	`json:"CAT"`
	TransactionId	string	`json:"transactionId"`
	Status			string	`json:"status"`
}

type Customer struct {
	NationalityId   string        `json:"nationalityId"`
    FamilyId        string        `json:"familyId"`
    Name            string        `json:"name"`
    Email           string        `json:"email"`
    PhoneNumber     string        `json:"phoneNumber"`
    QuotaRemaining  QuotaRemaining `json:"quotaRemaining"`
    CustomerTypes   []CustomerType `json:"customerTypes"`
    ChannelInject   string        `json:"channelInject"`
    IsAgreedTermsConditions bool   `json:"isAgreedTermsConditions"`
    IsCompleted    	bool           `json:"isCompleted"`
    IsSubsidi      bool           `json:"isSubsidi"`
}

type CustomerType struct {
    Name         	string 			`json:"name"`
    SourceTypeId 	int				`json:"sourceTypeId"`
    Status       	int    			`json:"status"`
    Verifications 	[]interface{} 	`json:"verifications"`
}

type QuotaRemaining struct {
    Individu int `json:"individu"`
    Family   int `json:"family"`
}

type Products struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}


type PengambilanItemSubsidi struct {
	Item           string `json:"item"`
	Quantitas      int    `json:"quantitas"`
	PotonganHarga  int    `json:"potongan_harga"`
}

type Subsidi struct {
	NIK                   string                `json:"nik"`
	IDValidation          string                `json:"IDValidation"`
	FamilyID              string                `json:"familyId"`
	Category              string                `json:"category"`
	SourceTypeID          int                   `json:"sourceTypeId"`
	Nama                  string                `json:"nama"`
	NoHandPhoneKPM        string                `json:"noHandPhoneKPM"`
	ChannelInject         string                `json:"channelInject"`
	PengambilanItemSubsidi []PengambilanItemSubsidi `json:"pengambilanItemSubsidi"`
}

type TransactionParam struct {
	Products      []Products `json:"products"`
	GeoTagging    string    `json:"geoTagging"`
	InputNominal  int       `json:"inputNominal"`
	Change        int       `json:"change"`
	PaymentType   string    `json:"paymentType"`
	Subsidi       Subsidi   `json:"subsidi"`
}

type TransactionSucces struct {
	TransactionId string `json:"transactionId"`
	TransactionUniqKey	string `json:"transactionUniqKey"`
}