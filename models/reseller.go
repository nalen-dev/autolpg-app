package models

type Reseller struct {
	RegistrationID           string  `json:"registrationId"`
	Name                     string  `json:"name"`
	Address                  string  `json:"address"`
	City                     string  `json:"city"`
	Province                 string  `json:"province"`
	Coordinate               string  `json:"coordinate"`
	StoreName                string  `json:"storeName"`
	StoreAddress             string  `json:"storeAddress"`
	PhoneNumber              string  `json:"phoneNumber"`
	TID                      string  `json:"tid"`
	MID                      *string `json:"mid"`
	SPBU                     string  `json:"spbu"`
	MerchantType             string  `json:"merchantType"`
	MIDMap                   string  `json:"midMap"`
	IsSubsidiProduct         bool    `json:"isSubsidiProduct"`
	StorePhoneNumber         string  `json:"storePhoneNumber"`
	Email                    string  `json:"email"`
	NationalityID            string  `json:"nationalityId"`
	DistrictName             string  `json:"ditrictName"`
	VillageName              string  `json:"villageName"`
	ZipCode                  string  `json:"zipcode"`
	Agen                     Agen    `json:"agen"`
	IsActiveMyptm            bool    `json:"isActiveMyptm"`
	Bank                     Bank    `json:"bank"`
	MyptmActivationStatus    *string `json:"myptmActivationStatus"`
	IsAvailableTransaction   bool    `json:"isAvailableTransaction"`
}

type Agen struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Bank struct {
	BankName     *string `json:"bankName"`
	AccountName  *string `json:"accountName"`
	AccountNumber *string `json:"accountNumber"`
}