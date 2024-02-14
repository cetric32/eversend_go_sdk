package eversendSdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Eversend struct
type Eversend struct {
	ClientId     string
	ClientSecret string
	BaseUrl      string
	Token        string
}

// NewEversend function to create a new Eversend instance
func NewEversendApp(clientId string, clientSecret string) *Eversend {
	return &Eversend{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		BaseUrl:      "https://api.eversend.co/v1/",
	}
}

func (e *Eversend) GenerateAuthToken() (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", e.BaseUrl+"auth/token", nil)

	if err != nil {
		return "", err
	}

	req.Header.Add("clientId", e.ClientId)
	req.Header.Add("clientSecret", e.ClientSecret)

	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	statusCode := resp.StatusCode
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	var responseData = map[string]interface{}{
		// "status":  200,
		"token":   "",
		"message": "",
	}

	err = json.Unmarshal(body, &responseData)

	if err != nil {
		return "", err
	}

	if statusCode != 200 {
		return "", errors.New(responseData["message"].(string))
	}

	// fmt.Println("Response Data:", responseData)

	token := responseData["token"].(string)

	e.Token = token

	return token, nil
}

func (e *Eversend) GetWallets() ([]interface{}, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", e.BaseUrl+"wallets", nil)

	if err != nil {
		return nil, err
	}

	token := e.Token

	if token == "" {
		token, err = e.GenerateAuthToken()

		if err != nil {
			return nil, err
		}
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	statusCode := resp.StatusCode
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var responseData = map[string]interface{}{
		// "status":  200,
		"data":    []interface{}{},
		"message": "",
	}

	err = json.Unmarshal(body, &responseData)

	if err != nil {
		return nil, err
	}

	if statusCode != 200 {
		return nil, errors.New(responseData["message"].(string))
	}

	// fmt.Println("Response Data:", responseData)

	data := responseData["data"].([]interface{})

	return data, nil
}

func (e *Eversend) CreateExchangeQuotation(from string, amount float64, to string) (map[string]interface{}, error) {
	client := &http.Client{}

	reqBody := []byte(fmt.Sprintf(`{"from": "%s", "amount": %f, "to": "%s"}`, from, amount, to))

	req, err := http.NewRequest("POST", e.BaseUrl+"exchanges/quotation", bytes.NewBuffer(reqBody))

	if err != nil {
		return nil, err
	}

	token := e.Token

	if token == "" {
		token, err = e.GenerateAuthToken()

		if err != nil {
			return nil, err
		}
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	statusCode := resp.StatusCode
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var responseData = map[string]interface{}{
		// "status":  200,
		"data":    map[string]interface{}{},
		"message": "",
	}

	err = json.Unmarshal(body, &responseData)

	if err != nil {
		return nil, err
	}

	if statusCode != 200 {
		return nil, errors.New(responseData["message"].(string))
	}

	// fmt.Println("Response Data:", responseData)

	data := responseData["data"].(map[string]interface{})

	return data, nil
}
