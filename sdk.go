package eversendSdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/cetric32/GoHTTP"
)

var eversendClientId string
var eversendClientSecret string
var baseUrl string = "https://api.eversend.co/v1/"
var authToken string
var mutex = &sync.RWMutex{}
var authTokenExpires time.Time

// Eversend struct
type Eversend struct {
	// clientId     string
	// clientSecret string
	// baseUrl      string
	// authToken    string

	Crypto        Crypto
	Wallets       Wallet
	Exchange      Exchange
	Payouts       Payout
	Beneficiaries Beneficiary
}

type Crypto struct{}
type Wallet struct{}
type Exchange struct{}
type Payout struct{}
type Beneficiary struct{}

// NewEversend function to create a new Eversend instance
func NewEversendApp(clientId string, clientSecret string) *Eversend {
	mutex.Lock()
	defer mutex.Unlock()

	eversendClientId = clientId
	eversendClientSecret = clientSecret

	return &Eversend{}
}

func generateAuthToken() (string, error) {
	//current time now UTC
	currentTime := time.Now()

	mutex.RLock()
	clientId := eversendClientId
	clientSecret := eversendClientSecret
	currentToken := authToken
	currentTokenExpires := authTokenExpires
	mutex.RUnlock()

	if currentToken != "" && currentTokenExpires.After(currentTime) {
		fmt.Println("existing token", currentToken, "expires", currentTokenExpires, "currentTime", currentTime)
		return currentToken, nil
	}

	url := baseUrl + "auth/token"

	goHttp := GoHTTP.NewGoHTTP()

	goHttp.AddHeaders(map[string]string{
		"clientId":     clientId,
		"clientSecret": clientSecret,
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
	expires := responseData["expires"].(string)

	mutex.Lock()
	defer mutex.Unlock()
	authToken = token
	authTokenExpires, err = time.Parse(time.RFC3339, expires)

	if err != nil {
		fmt.Println("error parsing expires", expires, err)
	}

	fmt.Println("new token", authToken, "expires", authTokenExpires, "currentTime", currentTime)
	return token, nil
}

// List function to fetch your eversend wallets and their balances
func (e *Wallet) List() ([]interface{}, error) {
	token, err := generateAuthToken()

	if err != nil {
		return nil, err
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

// Find function to fetch a specific Wallet and its balance
// The walletCurrency is the currency of the Wallet you want to get e.g "UGX"
func (e *Wallet) Find(walletCurrency string) (map[string]interface{}, error) {
	token, err := generateAuthToken()

	if err != nil {
		return nil, err
	}

	url := baseUrl + "wallets/" + walletCurrency

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

	data := responseData["data"].(map[string]interface{})

	return data, nil
}

// Quotation function to create an exchange quotation. This is used to get the amount you will receive when you convert money from one currency to another.
// It also gives you the exchange token which is used to create an exchange transaction.
// The amount is the amount you want to convert.
// The from is the currency you want to convert from e.g "UGX".
// The to is the currency you want to convert to e.g "KES".
func (e *Exchange) Quotation(from string, amount float64, to string) (map[string]interface{}, error) {
	url := baseUrl + "exchanges/quotation"
	token, err := generateAuthToken()

	if err != nil {
		return nil, err
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

// Exchange function to create an exchange transaction. This is used to convert money from one currency to another.
// The exchange token is used to identify the transaction. The exchange token is got from the CreateExchangeQuotation function
func (e *Exchange) Exchange(exchangeToken string) (map[string]interface{}, error) {
	url := baseUrl + "exchanges"
	token, err := generateAuthToken()

	if err != nil {
		return nil, err
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
	token, err := generateAuthToken()

	if err != nil {
		return nil, err
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

// DeliveryCountries function to get delivery countries. This are the countries you can send money to currently
func (e *Payout) DeliveryCountries() ([]interface{}, error) {
	url := baseUrl + "payouts/countries"
	token, err := generateAuthToken()

	if err != nil {
		return nil, err
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

// DeliveryBanks function to get delivery banks. This are the banks you can send money to in a specific country.
// The countryCode is the Alpha-2 country code of the country you want to get the banks for.
func (e *Payout) DeliveryBanks(countryCode string) ([]interface{}, error) {
	url := baseUrl + "payouts/banks/" + countryCode
	token, err := generateAuthToken()

	if err != nil {
		return nil, err
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

// Quotation function to create a Payout quotation. This is used to get the amount you will get and fees when you send money to a specific country.
// The amountType can be "DESTINATION" or "SOURCE". If it is "SOURCE", the amount is the amount that you want to be send. If it is "DESTINATION", the amount is the amount you want to be received.
// The Default is "SOURCE".
// The transactionType can be "bank" or "momo".
func (e *Payout) Quotation(sourceWallet string, amount float64,
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
	token, err := generateAuthToken()

	if err != nil {
		return nil, err
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

// MomoPayout function to create a mobile money(momo) Payout transaction. This is used to send money to a mobile money account of the recipient.
func (e *Payout) MomoPayout(payoutToken string, phoneNumber string, firstName string, lastName string, countryCode string) (map[string]interface{}, error) {
	url := baseUrl + "payouts"
	token, err := generateAuthToken()

	if err != nil {
		return nil, err
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

// BankPayout function to create a bank Payout transaction. This is used to send money to a bank account of the recipient.
func (e *Payout) BankPayout(payoutToken string, phoneNumber string, firstName string, lastName string,
	countryCode string, bankName string, bankAccountName string, bankCode string, bankAccountNumber string) (map[string]interface{}, error) {
	url := baseUrl + "payouts"

	token, err := generateAuthToken()

	if err != nil {
		return nil, err
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

// Transaction function to get a transaction details.
// The transactionId is the id of the transaction you want to get details for.
func (e *Payout) Transaction(transactionId string) (map[string]interface{}, error) {
	url := baseUrl + "transactions/" + transactionId
	token, err := generateAuthToken()

	if err != nil {
		return nil, err
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
func (e *Beneficiary) CreateMomoBeneficiary(firstName string, lastname string, countryCode string, phoneNumber string) error {
	url := baseUrl + "beneficiaries"

	token, err := generateAuthToken()

	if err != nil {
		return err
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
func (e *Beneficiary) CreateBankBeneficiary(firstName string, lastname string, countryCode string, bankName string,
	bankAccountName string, bankCode string, bankAccountNumber string) error {
	url := baseUrl + "beneficiaries"

	token, err := generateAuthToken()

	if err != nil {
		return err
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

// List function to get a list of beneficiaries. This is used to get the beneficiaries you have saved.
func (e *Beneficiary) List() ([]interface{}, error) {
	url := baseUrl + "beneficiaries"
	token, err := generateAuthToken()

	if err != nil {
		return nil, err
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

// Find function to get a beneficiary details. This is used to get the details of a specific Beneficiary.
func (e *Beneficiary) Find(beneficiaryId string) (map[string]interface{}, error) {
	url := baseUrl + "beneficiaries/" + beneficiaryId
	token, err := generateAuthToken()

	if err != nil {
		return nil, err
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

// AssetChains function to get a list of asset chains. This is used to get the asset chains you can use to send money.
// The coin is the currency you want to get the asset chains for e.g "USDT".
func (e *Crypto) AssetChains(coin string) (map[string]interface{}, error) {
	url := baseUrl + "crypto/assets/" + coin

	token, err := generateAuthToken()

	if err != nil {
		return nil, err
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

// Addresses function to get a list of addresses. This is used to get the addresses you have saved.
func (e *Crypto) Addresses() (map[string]interface{}, error) {
	url := baseUrl + "crypto/addresses"

	token, err := generateAuthToken()

	if err != nil {
		return nil, err
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

// GetTransactions function to get a list of crypto transactions. This is used to get the transactions you have made.
func (e *Crypto) Transactions() (map[string]interface{}, error) {
	url := baseUrl + "crypto/transactions"

	token, err := generateAuthToken()

	if err != nil {
		return nil, err
	}

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

	data := responseData["data"].(map[string]interface{})

	return data, nil
}

// AddressTransactions function to get a list of transactions for a specific address. This is used to get the transactions for a specific address.
func (e *Crypto) AddressTransactions(cryptoCoinAddress string) (map[string]interface{}, error) {
	url := baseUrl + "crypto/addresses/" + cryptoCoinAddress + "/transactions"

	token, err := generateAuthToken()

	if err != nil {
		return nil, err
	}

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

	data := responseData["data"].(map[string]interface{})

	return data, nil
}

// CreateAddress function to create a crypto address. This is used to create a crypto address for a specific coin.
// The assetId is the id of the asset you want to create the address for. Valid asset from the GetAssetChains function.
// The ownerName is the name of the owner of the address.
// The destinationAddressDescription is the description of the address. Should be Client email or unique identifier.
func (e *Crypto) CreateAddress(assetId string, ownerName string, destinationAddressDescription string, purpose string) (map[string]interface{}, error) {
	url := baseUrl + "crypto/addresses"

	token, err := generateAuthToken()

	if err != nil {
		return nil, err
	}
	reqBody := []byte(fmt.Sprintf(`{"assetId": "%s", "ownerName": "%s", "destinationAddressDescription": "%s", "purpose": "%s"}`,
		assetId, ownerName, destinationAddressDescription, purpose))

	goHttp := GoHTTP.NewGoHTTP()

	goHttp.AddHeaders(map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	})

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
