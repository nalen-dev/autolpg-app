package repository

import (
	"autolpg-app/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)


type ResellerRepository interface {
	GetResellerData() models.GetResellerResponse
}

type resellerRepository struct {
	httpClient *http.Client
	token 		string
}

func NewResellerRepo(httpClient *http.Client, token string) ResellerRepository{
	return &resellerRepository{
			httpClient: httpClient, 
			token: token,
		}
}

func (u resellerRepository) GetResellerData() models.GetResellerResponse {
	var response models.GetResellerResponse

	req, err := http.NewRequest(http.MethodGet, "https://api-map.my-pertamina.id/general/v1/users/profile", nil)
	if err != nil {
		log.Fatal(err)
		return response
	}

	req.Header.Set("Authorization", u.token)
	req.Header.Set("Origin", "https://subsiditepatlpg.mypertamina.id")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	resp, err := u.httpClient.Do(req)

	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return response
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

    err = json.Unmarshal(responseBody, &response)
    if err != nil {
        fmt.Println("Error decoding JSON:", err)
        return response
    }
	return response
}