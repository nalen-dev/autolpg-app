package models

type Product struct {
	RegistrationID   string    `json:"registrationId"`
	StoreName        string    `json:"storeName"`
	ProductID        string    `json:"productId"`
	ProductName      string    `json:"productName"`
	StockAvailable   int       `json:"stockAvailable"`
	StockRedeem      int       `json:"stockRedeem"`
	Sold             int       `json:"sold"`
	Modal            int       `json:"modal"`
	Price            int       `json:"price"`
	ProductMinPrice  int       `json:"productMinPrice"`
	ProductMaxPrice  int       `json:"productMaxPrice"`
	Image            string    `json:"image"`
	StockDate        string    `json:"stockDate"`
	LastStock        int       `json:"lastStock"`
	LastStockDate    string    `json:"lastStockDate"`
	LastSyncAt       string 	`json:"lastSyncAt"`
}