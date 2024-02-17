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
	clientId     string
	clientSecret string
	baseUrl      string
	authToken    string
}

// NewEversend function to create a new Eversend instance
func NewEversendApp(clientId string, clientSecret string) *Eversend {
	return &Eversend{
		clientId:     clientId,
		clientSecret: clientSecret,
		baseUrl:      "https://api.eversend.co/v1/",
	}
}

func (e *Eversend) generateAuthToken() (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", e.baseUrl+"auth/token", nil)

	if err != nil {
		return "", err
	}

	req.Header.Add("clientId", e.clientId)
	req.Header.Add("clientSecret", e.clientSecret)

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

	e.authToken = token

	return token, nil
}

// GetWallets function to fetch your eversend wallets and their balances
func (e *Eversend) GetWallets() ([]interface{}, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", e.baseUrl+"wallets", nil)

	if err != nil {
		return nil, err
	}

	token := e.authToken

	if token == "" {
		token, err = e.generateAuthToken()

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

// CreateExchangeQuotation function to create an exchange quotation. This is used to get the amount you will receive when you convert money from one currency to another.
// It also gives you the exchange token which is used to create an exchange transaction.
// The amount is the amount you want to convert.
// The from is the currency you want to convert from e.g "UGX".
// The to is the currency you want to convert to e.g "KES".
func (e *Eversend) CreateExchangeQuotation(from string, amount float64, to string) (map[string]interface{}, error) {
	client := &http.Client{}

	reqBody := []byte(fmt.Sprintf(`{"from": "%s", "amount": %f, "to": "%s"}`, from, amount, to))

	req, err := http.NewRequest("POST", e.baseUrl+"exchanges/quotation", bytes.NewBuffer(reqBody))

	if err != nil {
		return nil, err
	}

	token := e.authToken

	if token == "" {
		token, err = e.generateAuthToken()

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

// CreateExchange function to create an exchange transaction. This is used to convert money from one currency to another.
// The exchange token is used to identify the transaction. The exchange token is got from the CreateExchangeQuotation function
func (e *Eversend) CreateExchange(exchangeToken string) (map[string]interface{}, error) {
	client := &http.Client{}

	reqBody := []byte(fmt.Sprintf(`{"token": "%s"}`, exchangeToken))

	req, err := http.NewRequest("POST", e.baseUrl+"exchanges", bytes.NewBuffer(reqBody))

	if err != nil {
		return nil, err
	}

	token := e.authToken

	if token == "" {
		token, err = e.generateAuthToken()

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

	fmt.Println("Response Data:", responseData)

	data := responseData["data"].(map[string]interface{})

	return data, nil
}

// AccountProfile function to get account profile details
func (e *Eversend) AccountProfile() (map[string]interface{}, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", e.baseUrl+"account", nil)

	if err != nil {
		return nil, err
	}

	token := e.authToken

	if token == "" {
		token, err = e.generateAuthToken()

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

// GetDeliveryCountries function to get delivery countries. This are the countries you can send money to currently
func (e *Eversend) GetDeliveryCountries() ([]interface{}, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", e.baseUrl+"payouts/countries", nil)

	if err != nil {
		return nil, err
	}

	token := e.authToken

	if token == "" {
		token, err = e.generateAuthToken()

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

	return data["countries"].([]interface{}), nil
}

// GetDeliveryBanks function to get delivery banks. This are the banks you can send money to in a specific country.
// The countryCode is the Alpha-2 country code of the country you want to get the banks for.
func (e *Eversend) GetDeliveryBanks(countryCode string) ([]interface{}, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", e.baseUrl+"payouts/banks/"+countryCode, nil)

	if err != nil {
		return nil, err
	}

	token := e.authToken

	if token == "" {
		token, err = e.generateAuthToken()

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

// CreatePayoutQuotation function to create a payout quotation. This is used to get the amount you will get and fees when you send money to a specific country.
// The amountType can be "DESTINATION" or "SOURCE". If it is "SOURCE", the amount is the amount that you want to be send. If it is "DESTINATION", the amount is the amount you want to be received.
// The Default is "SOURCE".
// The transactionType can be "bank" or "momo".
func (e *Eversend) CreatePayoutQuotation(sourceWallet string, amount float64,
	transactionType string,
	destinationCountry string,
	destinationCurrency string,
	amountType string) (map[string]interface{}, error) {
	if amountType == "" {
		amountType = "SOURCE"
	}

	if amount < 0 {
		return nil, errors.New("amount cannot be negative")
	}

	client := &http.Client{}

	reqBody := []byte(fmt.Sprintf(`{"sourceWallet": "%s", "amount": %f, "type": "%s", "destinationCountry": "%s", "destinationCurrency": "%s", "amountType": "%s"}`,
		sourceWallet, amount, transactionType, destinationCountry, destinationCurrency, amountType))

	req, err := http.NewRequest("POST", e.baseUrl+"payouts/quotation", bytes.NewBuffer(reqBody))

	if err != nil {
		return nil, err
	}

	token := e.authToken

	if token == "" {
		token, err = e.generateAuthToken()

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

	println("Response Status Code:", resp.StatusCode)

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

// CreatePayout function to create a mobile money(momo) payout transaction. This is used to send money to a specific country.
func (e *Eversend) CreateMomoPayout(payoutToken string, phoneNumber string, firstName string, lastName string, countryCode string) (map[string]interface{}, error) {
	client := &http.Client{}

	reqBody := []byte(fmt.Sprintf(`{"token": "%s", "phoneNumber": "%s", "firstName": "%s", "lastName": "%s", "country": "%s"}`,
		payoutToken, phoneNumber, firstName, lastName, countryCode))

	req, err := http.NewRequest("POST", e.baseUrl+"payouts", bytes.NewBuffer(reqBody))

	if err != nil {
		return nil, err
	}

	token := e.authToken

	if token == "" {
		token, err = e.generateAuthToken()

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

// CreateBankPayout function to create a bank payout transaction. This is used to send money to a specific country.
func (e *Eversend) CreateBankPayout(payoutToken string, phoneNumber string, firstName string, lastName string,
	countryCode string, bankName string, bankAccountName string, bankCode string, bankAccountNumber string) (map[string]interface{}, error) {
	client := &http.Client{}

	reqBody := []byte(fmt.Sprintf(`{"token": "%s", "phoneNumber": "%s", "firstName": "%s", "lastName": "%s", "country": "%s", "bankName": "%s", "bankAccountName": "%s", "bankCode": "%s", "bankAccountNumber": "%s"}`,
		payoutToken, phoneNumber, firstName, lastName, countryCode, bankName, bankAccountName, bankCode, bankAccountNumber))

	req, err := http.NewRequest("POST", e.baseUrl+"payouts", bytes.NewBuffer(reqBody))

	if err != nil {
		return nil, err
	}

	token := e.authToken

	if token == "" {
		token, err = e.generateAuthToken()

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
