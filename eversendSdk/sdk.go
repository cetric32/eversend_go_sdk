package eversendSdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cetric32/eversend_go_sdk/GoHTTP"
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
	url := e.baseUrl + "auth/token"

	goHttp := GoHTTP.NewGoHTTP()

	goHttp.AddHeaders(map[string]string{
		"clientId":     e.clientId,
		"clientSecret": e.clientSecret,
	})

	body, statusCode, err := goHttp.Get(url)

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

	token := responseData["token"].(string)

	e.authToken = token

	return token, nil
}

// GetWallets function to fetch your eversend wallets and their balances
func (e *Eversend) GetWallets() ([]interface{}, error) {
	token := e.authToken
	var err error

	if token == "" {
		token, err = e.generateAuthToken()

		if err != nil {
			return nil, err
		}
	}

	url := e.baseUrl + "wallets"

	goHttp := GoHTTP.NewGoHTTP()

	goHttp.AddHeaders(map[string]string{
		"Authorization": "Bearer " + token,
	})

	body, statusCode, err := goHttp.Get(url)

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

	data := responseData["data"].([]interface{})

	return data, nil
}

// CreateExchangeQuotation function to create an exchange quotation. This is used to get the amount you will receive when you convert money from one currency to another.
// It also gives you the exchange token which is used to create an exchange transaction.
// The amount is the amount you want to convert.
// The from is the currency you want to convert from e.g "UGX".
// The to is the currency you want to convert to e.g "KES".
func (e *Eversend) CreateExchangeQuotation(from string, amount float64, to string) (map[string]interface{}, error) {
	url := e.baseUrl + "exchanges/quotation"
	var err error
	token := e.authToken

	if token == "" {
		token, err = e.generateAuthToken()

		if err != nil {
			return nil, err
		}
	}

	goHttp := GoHTTP.NewGoHTTP()

	goHttp.AddHeaders(map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	})

	reqBody, err := json.Marshal(map[string]interface{}{
		"from":   from,
		"amount": amount,
		"to":     to,
	})

	if err != nil {
		return nil, err
	}

	body, statusCode, err := goHttp.Post(url, bytes.NewBuffer(reqBody))

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

	data := responseData["data"].(map[string]interface{})

	return data, nil
}

// CreateExchange function to create an exchange transaction. This is used to convert money from one currency to another.
// The exchange token is used to identify the transaction. The exchange token is got from the CreateExchangeQuotation function
func (e *Eversend) CreateExchange(exchangeToken string) (map[string]interface{}, error) {
	url := e.baseUrl + "exchanges"
	var err error

	token := e.authToken

	if token == "" {
		token, err = e.generateAuthToken()

		if err != nil {
			return nil, err
		}
	}

	goHttp := GoHTTP.NewGoHTTP()

	goHttp.AddHeaders(map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	})

	reqBody := []byte(fmt.Sprintf(`{"token": "%s"}`, exchangeToken))

	body, statusCode, err := goHttp.Post(url, bytes.NewBuffer(reqBody))

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

	data := responseData["data"].(map[string]interface{})

	return data, nil
}

// AccountProfile function to get account profile details
func (e *Eversend) AccountProfile() (map[string]interface{}, error) {
	url := e.baseUrl + "account"
	var err error

	token := e.authToken

	if token == "" {
		token, err = e.generateAuthToken()

		if err != nil {
			return nil, err
		}
	}

	goHttp := GoHTTP.NewGoHTTP()

	goHttp.AddHeaders(map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	})

	body, statusCode, err := goHttp.Get(url)

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

	data := responseData["data"].(map[string]interface{})

	return data, nil
}

// GetDeliveryCountries function to get delivery countries. This are the countries you can send money to currently
func (e *Eversend) GetDeliveryCountries() ([]interface{}, error) {
	url := e.baseUrl + "payouts/countries"
	var err error

	token := e.authToken

	if token == "" {
		token, err = e.generateAuthToken()

		if err != nil {
			return nil, err
		}
	}

	goHttp := GoHTTP.NewGoHTTP()

	goHttp.AddHeaders(map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	})

	body, statusCode, err := goHttp.Get(url)

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

	data := responseData["data"].(map[string]interface{})

	return data["countries"].([]interface{}), nil
}

// GetDeliveryBanks function to get delivery banks. This are the banks you can send money to in a specific country.
// The countryCode is the Alpha-2 country code of the country you want to get the banks for.
func (e *Eversend) GetDeliveryBanks(countryCode string) ([]interface{}, error) {
	url := e.baseUrl + "payouts/banks/" + countryCode
	var err error

	token := e.authToken

	if token == "" {
		token, err = e.generateAuthToken()

		if err != nil {
			return nil, err
		}
	}

	goHttp := GoHTTP.NewGoHTTP()

	goHttp.AddHeaders(map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	})

	body, statusCode, err := goHttp.Get(url)

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

	url := e.baseUrl + "payouts/quotation"
	var err error

	token := e.authToken

	if token == "" {
		token, err = e.generateAuthToken()

		if err != nil {
			return nil, err
		}
	}

	goHttp := GoHTTP.NewGoHTTP()

	goHttp.AddHeaders(map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	})

	reqBody := []byte(fmt.Sprintf(`{"sourceWallet": "%s", "amount": %f, "type": "%s", "destinationCountry": "%s", "destinationCurrency": "%s", "amountType": "%s"}`,
		sourceWallet, amount, transactionType, destinationCountry, destinationCurrency, amountType))

	body, statusCode, err := goHttp.Post(url, bytes.NewBuffer(reqBody))

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

	data := responseData["data"].(map[string]interface{})

	return data, nil
}

// CreatePayout function to create a mobile money(momo) payout transaction. This is used to send money to a mobile money account of the recipient.
func (e *Eversend) CreateMomoPayout(payoutToken string, phoneNumber string, firstName string, lastName string, countryCode string) (map[string]interface{}, error) {
	url := e.baseUrl + "payouts"
	var err error

	token := e.authToken

	if token == "" {
		token, err = e.generateAuthToken()

		if err != nil {
			return nil, err
		}
	}

	goHttp := GoHTTP.NewGoHTTP()

	goHttp.AddHeaders(map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	})

	reqBody := []byte(fmt.Sprintf(`{"token": "%s", "phoneNumber": "%s", "firstName": "%s", "lastName": "%s", "country": "%s"}`,
		payoutToken, phoneNumber, firstName, lastName, countryCode))

	body, statusCode, err := goHttp.Post(url, bytes.NewBuffer(reqBody))

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

	data := responseData["data"].(map[string]interface{})
	return data, nil
}

// CreateBankPayout function to create a bank payout transaction. This is used to send money to a bank account of the recipient.
func (e *Eversend) CreateBankPayout(payoutToken string, phoneNumber string, firstName string, lastName string,
	countryCode string, bankName string, bankAccountName string, bankCode string, bankAccountNumber string) (map[string]interface{}, error) {
	url := e.baseUrl + "payouts"

	var err error

	token := e.authToken

	if token == "" {
		token, err = e.generateAuthToken()

		if err != nil {
			return nil, err
		}
	}

	goHttp := GoHTTP.NewGoHTTP()

	goHttp.AddHeaders(map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	})

	reqBody := []byte(fmt.Sprintf(`{"token": "%s", "phoneNumber": "%s", "firstName": "%s", "lastName": "%s", "country": "%s", "bankName": "%s", "bankAccountName": "%s", "bankCode": "%s", "bankAccountNumber": "%s"}`,
		payoutToken, phoneNumber, firstName, lastName, countryCode, bankName, bankAccountName, bankCode, bankAccountNumber))

	body, statusCode, err := goHttp.Post(url, bytes.NewBuffer(reqBody))

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

	data := responseData["data"].(map[string]interface{})
	return data, nil
}

// GetTransaction function to get a transaction details.
// The transactionId is the id of the transaction you want to get details for.
func (e *Eversend) GetTransaction(transactionId string) (map[string]interface{}, error) {
	url := e.baseUrl + "transactions/" + transactionId
	var err error

	token := e.authToken

	if token == "" {
		token, err = e.generateAuthToken()

		if err != nil {
			return nil, err
		}
	}

	goHttp := GoHTTP.NewGoHTTP()

	goHttp.AddHeaders(map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	})

	body, statusCode, err := goHttp.Get(url)

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

	data := responseData["data"].(map[string]interface{})
	return data, nil
}
