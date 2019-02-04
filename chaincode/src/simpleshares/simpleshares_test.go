package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
)

func Test_Init(t *testing.T) {
	fmt.Println("Testing Chaincode Initialisation...")
	simpleSharesCC := new(SimpleShares)
	mockStub := shim.NewMockStub("mockstub", simpleSharesCC)
	txID := "mockTxID"

	mockStub.MockTransactionStart(txID)
	response := simpleSharesCC.Init(mockStub)
	mockStub.MockTransactionEnd(txID)
	if s := response.GetStatus(); s != 200 {
		fmt.Println("Init test failed")
		t.FailNow()
	}
}

func Test_Issue(t *testing.T) {
	fmt.Println("Testing Share Issuance...")
	simpleSharesCC := new(SimpleShares)
	mockStub := shim.NewMockStub("mockstub", simpleSharesCC)
	txID := "mockTxID"

	// Initialise chaincode
	mockStub.MockTransactionStart(txID)
	_ = simpleSharesCC.Init(mockStub)

	// Issue new shares and append to DB
	args := []string{"issue", "NewOrg", "1000", "0xe7cf944311eabff15b1b091422a2ecada1dd053d"}
	mockStub.MockTransactionStart(txID)
	response := simpleSharesCC.issue(mockStub, args)
	mockStub.MockTransactionEnd(txID)

	assert.Equal(t, response.GetStatus(), int32(200))

	// Query DB for new state
	args = []string{"query", "shares"}
	mockStub.MockTransactionStart(txID)
	response = simpleSharesCC.query(mockStub, args)

	var listings Shares
	err := json.Unmarshal(response.GetPayload(), &listings)
	if err != nil {
		fmt.Println("error:", err)
	}

	assert.Equal(t, uint(1000), listings.Organisation[0].TotalSupply)
	assert.Equal(t, common.HexToAddress("0xe7cf944311eabff15b1b091422a2ecada1dd053d"), listings.Organisation[0].Shareholders[0].Shareholder, "Gosbank")
	assert.Equal(t, uint(1000), listings.Organisation[0].Shareholders[0].Amount)

}

func Test_Transfer(t *testing.T) {
	fmt.Println("Testing Title Transfer...")
	simpleSharesCC := new(SimpleShares)
	mockStub := shim.NewMockStub("mockstub", simpleSharesCC)
	txID := "mockTxID"

	mockStub.MockTransactionStart(txID)
	_ = simpleSharesCC.Init(mockStub)

	// initial share issueing required as DB is empty
	args := []string{"issue", "NewOrg", "1000", "0xe7cf944311eabff15b1b091422a2ecada1dd053d"}
	mockStub.MockTransactionStart(txID)
	_ = simpleSharesCC.issue(mockStub, args)
	mockStub.MockTransactionEnd(txID)

	// place orders to buy/sell equities for each counterparty
	mockStub.MockTransactionStart(txID)
	args = []string{"order", "NewOrg", "Buy", "5", "100", "0x9ecd4a8ca1560c4bb92ca9ebfa5eab448048db93", "uniqueRef"}
	_ = simpleSharesCC.order(mockStub, args)
	args = []string{"order", "NewOrg", "Sell", "5", "100", "0xe7cf944311eabff15b1b091422a2ecada1dd053d", "uniqueRef"}
	_ = simpleSharesCC.order(mockStub, args)
	mockStub.MockTransactionEnd(txID)

	args = []string{"transfer", "uniqueRef"}
	mockStub.MockTransactionStart(txID)
	response := simpleSharesCC.transfer(mockStub, args)
	mockStub.MockTransactionEnd(txID)

	assert.Equal(t, int32(200), response.GetStatus())

	args = []string{"query", "shares"}
	response = simpleSharesCC.query(mockStub, args)

	var shares Shares
	err := json.Unmarshal(response.GetPayload(), &shares)
	if err != nil {
		fmt.Println("error:", err)
	}

	assert.Equal(t, uint(1000), shares.Organisation[0].TotalSupply)
	assert.Equal(t, common.HexToAddress("0xe7cf944311eabff15b1b091422a2ecada1dd053d"), shares.Organisation[0].Shareholders[0].Shareholder)
	assert.Equal(t, uint(500), shares.Organisation[0].Shareholders[0].Amount)
	assert.Equal(t, common.HexToAddress("0x9ecd4a8ca1560c4bb92ca9ebfa5eab448048db93"), shares.Organisation[0].Shareholders[1].Shareholder)
	assert.Equal(t, uint(500), shares.Organisation[0].Shareholders[1].Amount)

	args = []string{"query", "trades"}
	response = simpleSharesCC.query(mockStub, args)

	var trades Trades
	err = json.Unmarshal(response.GetPayload(), &trades)
	if err != nil {
		fmt.Println("error:", err)
	}

	assert.Equal(t, "uniqueRef", trades.Trade[0].Ref)

}

// func Test_InsufficientBalance(t *testing.T) {
// 	fmt.Println("Testing Transfer - Sender Insufficient Balance...")
// 	simpleSharesCC := new(SimpleShares)
// 	mockStub := shim.NewMockStub("mockstub", simpleSharesCC)
// 	txID := "mockTxID"

// 	mockStub.MockTransactionStart(txID)
// 	_ = simpleSharesCC.Init(mockStub)

// 	// initial share issueing required as DB is empty
// 	args := []string{"issue", "NewOrg", "1000", "0xe7cf944311eabff15b1b091422a2ecada1dd053d"}
// 	mockStub.MockTransactionStart(txID)
// 	_ = simpleSharesCC.issue(mockStub, args)
// 	mockStub.MockTransactionEnd(txID)

// 	args = []string{"transfer", "NewOrg", "1500", "0xe7cf944311eabff15b1b091422a2ecada1dd053d", "0x9ecd4a8ca1560c4bb92ca9ebfa5eab448048db93"}
// 	mockStub.MockTransactionStart(txID)
// 	response := simpleSharesCC.transfer(mockStub, args)
// 	mockStub.MockTransactionEnd(txID)

// 	assert.Equal(t, int32(500), response.GetStatus())
// 	assert.Equal(t, "Insufficient Balance", response.GetMessage())
// }

// func Test_UnknownSender(t *testing.T) {
// 	fmt.Println("Testing Transfer - Unknown Sender...")
// 	simpleSharesCC := new(SimpleShares)
// 	mockStub := shim.NewMockStub("mockstub", simpleSharesCC)
// 	txID := "mockTxID"

// 	mockStub.MockTransactionStart(txID)
// 	_ = simpleSharesCC.Init(mockStub)

// 	// initial share issueing required as DB is empty
// 	args := []string{"issue", "NewOrg", "1000", "0xe7cf944311eabff15b1b091422a2ecada1dd053d"}
// 	mockStub.MockTransactionStart(txID)
// 	_ = simpleSharesCC.issue(mockStub, args)
// 	mockStub.MockTransactionEnd(txID)

// 	args = []string{"transfer", "NewOrg", "500", "0xe6cf944311eabff15b1b091422a2ecada1dd053d", "0x9ecd4a8ca1560c4bb92ca9ebfa5eab448048db93"}
// 	mockStub.MockTransactionStart(txID)
// 	response := simpleSharesCC.transfer(mockStub, args)
// 	mockStub.MockTransactionEnd(txID)

// 	assert.Equal(t, int32(500), response.GetStatus())
// 	assert.Equal(t, "Insufficient Balance", response.GetMessage())
// }

func Test_PlaceOrder(t *testing.T) {
	fmt.Println("Testing Order Placing...")
	simpleSharesCC := new(SimpleShares)
	mockStub := shim.NewMockStub("mockstub", simpleSharesCC)
	txID := "mockTxID"

	mockStub.MockTransactionStart(txID)
	_ = simpleSharesCC.Init(mockStub)

	args := []string{"order", "NewOrg", "Buy", "100", "100", "0xe7cf944311eabff15b1b091422a2ecada1dd053d", "uniqueRef"}
	mockStub.MockTransactionStart(txID)
	response := simpleSharesCC.order(mockStub, args)
	mockStub.MockTransactionEnd(txID)

	args = []string{"query", "orders"}
	mockStub.MockTransactionStart(txID)
	response = simpleSharesCC.query(mockStub, args)

	assert.Equal(t, int32(200), response.GetStatus())

	var orderBook OrderStruct
	err := json.Unmarshal(response.GetPayload(), &orderBook)
	if err != nil {
		fmt.Println("error:", err)
	}

	assert.Equal(t, "NewOrg", orderBook.Open[0].Organisation)
	assert.Equal(t, "Buy", orderBook.Open[0].BuySell)
	assert.Equal(t, uint(100), orderBook.Open[0].Amount)
	assert.Equal(t, uint(100), orderBook.Open[0].Price)
	assert.Equal(t, common.HexToAddress("0xe7cf944311eabff15b1b091422a2ecada1dd053d"), orderBook.Open[0].Address)
}
func Test_MatchOrder(t *testing.T) {
	fmt.Println("Testing Order Matching...")
	simpleSharesCC := new(SimpleShares)
	mockStub := shim.NewMockStub("mockstub", simpleSharesCC)
	txID := "mockTxID"

	mockStub.MockTransactionStart(txID)
	_ = simpleSharesCC.Init(mockStub)

	args := []string{"order", "NewOrg", "Buy", "100", "100", "0xe7cf944311eabff15b1b091422a2ecada1dd053d", "uniqueRef"}
	mockStub.MockTransactionStart(txID)
	_ = simpleSharesCC.order(mockStub, args)
	mockStub.MockTransactionEnd(txID)

	args = []string{"order", "NewOrg", "Sell", "100", "100", "0x9ecd4a8ca1560c4bb92ca9ebfa5eab448048db93", "uniqueRef"}
	mockStub.MockTransactionStart(txID)
	_ = simpleSharesCC.order(mockStub, args)
	mockStub.MockTransactionEnd(txID)

	args = []string{"query", "orders"}
	mockStub.MockTransactionStart(txID)
	response := simpleSharesCC.query(mockStub, args)

	var orderBook OrderStruct
	err := json.Unmarshal(response.GetPayload(), &orderBook)
	if err != nil {
		fmt.Println("error:", err)
	}

	assert.Equal(t, "NewOrg", orderBook.Matched[0].Organisation)
	assert.Equal(t, uint(100), orderBook.Matched[0].Amount)
	assert.Equal(t, uint(100), orderBook.Matched[0].Price)
	assert.Equal(t, common.HexToAddress("0x9ecd4a8ca1560c4bb92ca9ebfa5eab448048db93"), orderBook.Matched[0].Send)
	assert.Equal(t, common.HexToAddress("0xe7cf944311eabff15b1b091422a2ecada1dd053d"), orderBook.Matched[0].Recv)
	assert.Equal(t, "uniqueRef", orderBook.Matched[0].Ref)
}
