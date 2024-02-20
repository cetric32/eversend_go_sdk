package eversendSdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/cetric32/eversend_go_sdk/GoHTTP"
)

var eversendClientId string
var eversendClientSecret string
var baseUrl string = "https://api.eversend.co/v1/"
var authToken string
var mutex = &sync.Mutex{}

// Eversend struct
type Eversend struct {
	// clientId     string
	// clientSecret string
	// baseUrl      string
	// authToken    string

	Crypto crypto
}

type crypto struct{}

// NewEversend function to create a new Eversend instance
func NewEversendApp(clientId string, clientSecret string) *Eversend {
	mutex.Lock()
	defer mutex.Unlock()

	eversendClientId = clientId
	eversendClientSecret = clientSecret

	return &Eversend{}
}

func generateAuthToken() (string, error) {
	url := baseUrl + "auth/token"

	goHttp := GoHTTP.NewGoHTTP()

	goHttp.AddHeaders(map[string]string{
		"clientId":     eversendClientId,
		"clientSecret": eversendClientSecret,
	})

	body, statusCode, err := goHttp.Get(url)

	if err != nil {
		return "", err
	}

	var responseData = map[string]interface{}{}

	err = json.Unmarshal(body, &responseData)

	if err != nil {
		return "", err
	}

	if statusCode != 200 {
		return "", errors.New(responseData["message"].(string))
	}

	token := responseData["token"].(string)

	mutex.Lock()
	defer mutex.Unlock()
	authToken = token

	return token, nil
}

// GetWallets function to fetch your eversend wallets and their balances
func (e *Eversend) GetWallets() ([]interface{}, error) {
	token := authToken
	var err error

	if token == "" {
		token, err = generateAuthToken()

		if err != nil {
			return nil, err
		}
	}

	url := baseUrl + "wallets"

	goHttp := GoHTTP.NewGoHTTP()

	goHttp.AddHeaders(map[string]string{
		"Authorization": "Bearer " + token,
	})

	body, statusCode, err := goHttp.Get(url)

	if err != nil {
		return nil, err
	}

	var responseData = map[string]interface{}{}

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
	url := baseUrl + "exchanges/quotation"
	var err error
	token := authToken

	if token == "" {
		token, err = generateAuthToken()

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

	var responseData = map[string]interface{}{}

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
	url := baseUrl + "exchanges"
	var err error

	token := authToken

	if token == "" {
		token, err = generateAuthToken()

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

	var responseData = map[string]interface{}{}

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
	url := baseUrl + "account"
	var err error

	token := authToken

	if token == "" {
		token, err = generateAuthToken()

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

	var responseData = map[string]interface{}{}

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
	url := baseUrl + "payouts/countries"
	var err error

	token := authToken

	if token == "" {
		token, err = generateAuthToken()

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

	var responseData = map[string]interface{}{}

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
	url := baseUrl + "payouts/banks/" + countryCode
	var err error

	token := authToken

	if token == "" {
		token, err = generateAuthToken()

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

	var responseData = map[string]interface{}{}

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

	url := baseUrl + "payouts/quotation"
	var err error

	token := authToken

	if token == "" {
		token, err = generateAuthToken()

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

	var responseData = map[string]interface{}{}

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
	url := baseUrl + "payouts"
	var err error

	token := authToken

	if token == "" {
		token, err = generateAuthToken()

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

	var responseData = map[string]interface{}{}

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
	url := baseUrl + "payouts"

	var err error

	token := authToken

	if token == "" {
		token, err = generateAuthToken()

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

	var responseData = map[string]interface{}{}

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
	url := baseUrl + "transactions/" + transactionId
	var err error

	token := authToken

	if token == "" {
		token, err = generateAuthToken()

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

	var responseData = map[string]interface{}{}

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

// CreateMomoBeneficiary function to create a mobile money beneficiary. This is used to save a mobile money account for future use.
// countryCode is the Alpha-2 country code of the country e.g "UG".
func (e *Eversend) CreateMomoBeneficiary(firstName string, lastname string, countryCode string, phoneNumber string) error {
	url := baseUrl + "beneficiaries"

	var err error

	token := authToken

	if token == "" {
		token, err = generateAuthToken()

		if err != nil {
			return err
		}
	}

	goHttp := GoHTTP.NewGoHTTP()

	goHttp.AddHeaders(map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	})

	reqBody := []byte(fmt.Sprintf(`{"firstName": "%s", "lastName": "%s", "country": "%s", "phoneNumber": "%s","isBank": false,"isMomo": true}`,
		firstName, lastname, countryCode, phoneNumber))

	respBody, statusCode, err := goHttp.Post(url, bytes.NewBuffer(reqBody))

	if err != nil {
		return err
	}

	var responseData = map[string]interface{}{}

	err = json.Unmarshal(respBody, &responseData)

	if err != nil {
		return err
	}

	if statusCode != 200 {
		return errors.New(responseData["message"].(string))
	}

	return nil
}

// CreateBankBeneficiary function to create a bank beneficiary. This is used to save a bank account for future use.
// bankCode is got from the GetDeliveryBanks function.
// countryCode is the Alpha-2 country code of the country e.g "UG".
func (e *Eversend) CreateBankBeneficiary(firstName string, lastname string, countryCode string, bankName string,
	bankAccountName string, bankCode string, bankAccountNumber string) error {
	url := baseUrl + "beneficiaries"

	var err error

	token := authToken

	if token == "" {
		token, err = generateAuthToken()

		if err != nil {
			return err
		}
	}

	goHttp := GoHTTP.NewGoHTTP()

	goHttp.AddHeaders(map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	})

	reqBody := []byte(fmt.Sprintf(`{"firstName": "%s", "lastName": "%s", "country": "%s", "bankName": "%s", "bankAccountName": "%s", "bankCode": "%s", "bankAccountNumber": "%s","isBank": true,"isMomo": true}`,
		firstName, lastname, countryCode, bankName, bankAccountName, bankCode, bankAccountNumber))

	respBody, statusCode, err := goHttp.Post(url, bytes.NewBuffer(reqBody))

	if err != nil {
		return err
	}

	var responseData = map[string]interface{}{}

	err = json.Unmarshal(respBody, &responseData)

	if err != nil {
		return err
	}

	if statusCode != 200 {
		return errors.New(responseData["message"].(string))
	}

	return nil
}

// GetBeneficiaries function to get a list of beneficiaries. This is used to get the beneficiaries you have saved.
func (e *Eversend) GetBeneficiaries() ([]interface{}, error) {
	url := baseUrl + "beneficiaries"
	var err error

	token := authToken

	if token == "" {
		token, err = generateAuthToken()

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

	var responseData = map[string]interface{}{}

	err = json.Unmarshal(body, &responseData)

	if err != nil {
		return nil, err
	}

	if statusCode != 200 {
		return nil, errors.New(responseData["message"].(string))
	}

	data := responseData["data"].(map[string]interface{})

	beneficiaries := data["beneficiaries"].([]interface{})

	return beneficiaries, nil
}

func (e *Eversend) GetBeneficiary(beneficiaryId string) (map[string]interface{}, error) {
	url := baseUrl + "beneficiaries/" + beneficiaryId
	var err error

	token := authToken

	if token == "" {
		token, err = generateAuthToken()

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

	var responseData = map[string]interface{}{}

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

// GetAssetChains function to get a list of asset chains. This is used to get the asset chains you can use to send money.
// The coin is the currency you want to get the asset chains for e.g "USDT".
func (e *crypto) GetAssetChains(coin string) (map[string]interface{}, error) {
	url := baseUrl + "crypto/assets/" + coin

	var err error

	token := authToken

	if token == "" {
		token, err = generateAuthToken()

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

	var responseData = map[string]interface{}{}

	err = json.Unmarshal(body, &responseData)

	if err != nil {
		return nil, err
	}

	if statusCode != 200 {
		return nil, errors.New(responseData["message"].(string))
	}

	return responseData, nil
}

// GetAddresses function to get a list of addresses. This is used to get the addresses you have saved.
func (e *crypto) GetAddresses() (map[string]interface{}, error) {
	url := baseUrl + "crypto/addresses"

	var err error

	token := authToken

	if token == "" {
		token, err = generateAuthToken()

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

	var responseData = map[string]interface{}{}

	err = json.Unmarshal(body, &responseData)

	if err != nil {
		return nil, err
	}

	if statusCode != 200 {
		return nil, errors.New(responseData["message"].(string))
	}

	return responseData, nil
}
