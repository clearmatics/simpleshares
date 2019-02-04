package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/clearmatics/autonity/rlp"
	"github.com/ethereum/go-ethereum/common"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleShares implementation of Chaincode
type SimpleShares struct {
}

// Shares struct that logs the shares held by parties
type Shares struct {
	Organisation []Organisation `json:"organisation"`
}

type Organisation struct {
	Name         string         `json:"name"`
	TotalSupply  uint           `json:"totalsupply"`
	Shareholders []Shareholders `json:"shareholders"`
}

type Shareholders struct {
	Shareholder common.Address `json:"shareholder"`
	Amount      uint           `json:"amount"`
}

type OrderStruct struct {
	Open    []OpenOrder
	Matched []MatchedOrder
}

type OpenOrder struct {
	Organisation string
	BuySell      string
	Amount       uint
	Price        uint
	Address      common.Address
	Ref          string
}

type MatchedOrder struct {
	Organisation string
	Send         common.Address
	Recv         common.Address
	Amount       uint
	Price        uint
	Ref          string
}

type Trades struct {
	Trade []TradeStruct
}

type TradeStruct struct {
	Ref     string
	Details MatchedOrder
}

// EmptyAddress ugly helper to allow us to see if an address is empty
var EmptyAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")

// Init of the chaincode
// This function is called only one when the chaincode is instantiated.
// So the goal is to prepare the ledger to handle future requests.
func (t *SimpleShares) Init(stub shim.ChaincodeStubInterface) pb.Response {
	// Create Share listings
	listings := make([]interface{}, 0)
	shareConfig, err := rlp.EncodeToBytes(listings)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Create order book
	orderbook := make([]interface{}, 0)
	ordersConfig, err := rlp.EncodeToBytes(orderbook)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Create executed trade list
	trades := make([]interface{}, 0)
	tradesConfig, err := rlp.EncodeToBytes(trades)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState("shares", shareConfig)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState("orders", ordersConfig)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState("trades", tradesConfig)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return a successful message
	return shim.Success(nil)
}

// Invoke
// All future requests named invoke will arrive here.
func (t *SimpleShares) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	args := stub.GetStringArgs()
	if len(args) == 0 {
		return shim.Error("Function not provided")
	}

	// Check whether the number of arguments is sufficient
	if len(args) < 1 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// In order to manage multiple type of request, we will check the first argument.
	// Here we have one possible argument: query (every query request will read in the ledger without modification)
	if args[0] == "query" {
		return t.query(stub, args)
	}

	// input argument to issue a new share listing
	if args[0] == "issue" {
		return t.issue(stub, args)
	}

	// transfer funds
	if args[0] == "transfer" {
		return t.transfer(stub, args)
	}

	// place a new order
	if args[0] == "order" {
		return t.order(stub, args)
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown action, check the first argument")
}

// query
// usage: query [key]
func (t *SimpleShares) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Check whether the number of arguments is sufficient
	if len(args) != 2 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// Retrieve state of orderbook or share holdings
	state, err := stub.GetState(args[1])
	if err != nil {
		return shim.Error("Failed to get state")
	}

	if args[1] == "shares" {
		var listings Shares
		err = rlp.DecodeBytes(state, &listings)

		// Return bytes output of listings
		output, err := marshallState(listings)
		if err != nil {
			return shim.Error("Failed to marshall state JSON")
		}
		// Return this value in response
		return shim.Success(output)

	} else if args[1] == "orders" {
		var orderBook OrderStruct
		err = rlp.DecodeBytes(state, &orderBook)

		// Return bytes output of orderbook
		output, err := marshallOrderBook(orderBook)
		if err != nil {
			return shim.Error("Failed to marshall state JSON")
		}
		// Return this value in response
		return shim.Success(output)

	} else if args[1] == "trades" {
		var trades Trades
		err = rlp.DecodeBytes(state, &trades)

		// Return bytes output of trade listings
		output, err := marshallTrades(trades)
		if err != nil {
			return shim.Error("Failed to marshall state JSON")
		}
		// Return this value in response
		return shim.Success(output)

	} else {
		return shim.Error("Incorrect input argument")
	}

}

// issue
// usage: issue [shares] [amount] [recv]
func (t *SimpleShares) issue(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// Retrieve current state
	state, err := stub.GetState("shares")
	if err != nil {
		return shim.Error("Failed to get state")
	}

	var orderBook OrderStruct
	err = rlp.DecodeBytes(state, &orderBook)

	var listings Shares
	err = rlp.DecodeBytes(state, &listings)

	// retrieve number of shares to be issue
	numShares, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Failed convert balance")
	}

	// Convert string to address
	addr := common.HexToAddress(args[3])

	var newShareholder []Shareholders
	var newOrganisation []Organisation

	// create new shareholder struct and append to organisation
	newShareholder = append(newShareholder, Shareholders{Shareholder: addr, Amount: uint(numShares)})
	newOrganisation = append(newOrganisation, Organisation{Name: args[1], TotalSupply: uint(numShares), Shareholders: newShareholder})
	listings.Organisation = append(listings.Organisation, Organisation{Name: args[1], TotalSupply: uint(numShares), Shareholders: newShareholder})

	// Create new state and write to DB
	newState, err := rlp.EncodeToBytes(listings)
	if err != nil {
		return shim.Error("Failed to RLP encode state")
	}

	err = stub.PutState("shares", newState)
	if err != nil {
		return shim.Error("Failed to update state of simpleshares")
	}

	// Notify listeners that an event "issueEvent" have been executed
	err = stub.SetEvent("issueEvent", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// transfer
// transfers [ref]
func (t *SimpleShares) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// Retrieve orderbook state
	orderBook, err := retrieveOrderBook(stub)
	if err != nil {
		return shim.Error("Failed to get orderbook")
	}

	// Find trade details held with unique reference
	var orderDetails MatchedOrder
	for i := range orderBook.Matched {
		if orderBook.Matched[i].Ref == args[1] {
			orderDetails = orderBook.Matched[i]
		}
	}

	// Retrieve current state
	listings, err := retrieveShareListings(stub)
	if err != nil {
		return shim.Error("Failed to get share listings")
	}

	// Create new structs for sender and receiver
	var newOrganisation Organisation
	var send Shareholders
	var recv Shareholders

	// Retrieve infomation for specific issued shares
	for i := range listings.Organisation {
		if listings.Organisation[i].Name == orderDetails.Organisation {
			newOrganisation = listings.Organisation[i]
		}
	}
	// Check sender balance is sufficient
	if newOrganisation.Name == "" {
		return shim.Error("Unknown Organisation")
	}

	// Retrieve sender balance
	for i := range newOrganisation.Shareholders {
		if newOrganisation.Shareholders[i].Shareholder == orderDetails.Send {
			send = newOrganisation.Shareholders[i]
		}
		if newOrganisation.Shareholders[i].Shareholder == orderDetails.Recv {
			recv = newOrganisation.Shareholders[i]
		}
	}

	// Check sender balance is sufficient
	amount := orderDetails.Amount * orderDetails.Price
	if send.Amount < amount {
		return shim.Error("Sender Has Insufficient Balance")
	}

	// ugly way to check if the recv didn't previously exist but works for now
	if recv.Shareholder.String() == EmptyAddress.String() {
		recv.Shareholder = orderDetails.Recv
		newOrganisation.Shareholders = append(newOrganisation.Shareholders, recv)
	}

	// Transfer funds between sender and recipient
	send.Amount = send.Amount - uint(amount)
	recv.Amount = recv.Amount + uint(amount)

	// Update holdings of individual shareholders
	for i := range newOrganisation.Shareholders {
		if newOrganisation.Shareholders[i].Shareholder == send.Shareholder {
			newOrganisation.Shareholders[i] = send
		}
		if newOrganisation.Shareholders[i].Shareholder == recv.Shareholder {
			newOrganisation.Shareholders[i] = recv
		}
	}

	// Rewrite new state of share listings of organisation which has had equities traded
	for i := range listings.Organisation {
		if listings.Organisation[i].Name == orderBook.Matched[i].Organisation {
			listings.Organisation[i] = newOrganisation
		}
	}

	// Create new state and write to DB
	newState, err := rlp.EncodeToBytes(listings)
	if err != nil {
		return shim.Error("Failed to RLP encode state")
	}

	// Update shares state
	err = stub.PutState("shares", newState)
	if err != nil {
		return shim.Error("Failed to update state of simpleshares")
	}

	// Mark trade completed
	err = t.trade(stub, orderDetails)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Notify listeners that an event has been executed
	err = stub.SetEvent("transferEvent", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)

}

// Trade
func (t *SimpleShares) trade(stub shim.ChaincodeStubInterface, order MatchedOrder) error {
	// Retrieve trades state
	state, err := stub.GetState("trades")
	if err != nil {
		return errors.New("Failed to get state")
	}

	var trades Trades
	err = rlp.DecodeBytes(state, &trades)

	// create new executed trade, allows to find through reference
	var newTrade TradeStruct
	newTrade.Ref = order.Ref
	newTrade.Details = order

	// Search through trades check there is not a trade already with this reference
	// else append the trade to executed trade ledger
	for i := range trades.Trade {
		if trades.Trade[i].Ref == order.Ref {
			return errors.New("Trade Already Exists")
		}
	}

	trades.Trade = append(trades.Trade, newTrade)
	// Create new state and write to DB
	newState, err := rlp.EncodeToBytes(trades)
	if err != nil {
		return errors.New("Failed to RLP encode state")
	}

	fmt.Printf("%x\n", newState)

	// Update shares state
	err = stub.PutState("trades", newState)
	if err != nil {
		fmt.Println(err)
		return errors.New("Failed to update state of simpleshares")
	}

	return nil
}

// order
// usage: [organisation] [buy/sell] [amount] [price] [sender] [reference]
func (t *SimpleShares) order(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Verify enough arguments have been passed, could check formatting in future
	if len(args) != 7 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// Retrieve current state
	state, err := stub.GetState("orders")
	if err != nil {
		return shim.Error("Failed to get state")
	}

	var orderBook OrderStruct
	err = rlp.DecodeBytes(state, &orderBook)

	// Create new order from arguments
	newOrder, err := newOrder(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Attempt to match order
	matchedOrder := matchOrder(orderBook, newOrder)

	// If order is not matched append to open orders else remove from open orders and append to matched orders
	if matchedOrder.Organisation == "" {
		orderBook.Open = append(orderBook.Open, newOrder)
	} else {
		orderBook = orderMatched(orderBook, matchedOrder, newOrder)
	}

	// Create new state and write to DB
	newState, err := rlp.EncodeToBytes(orderBook)
	if err != nil {
		return shim.Error("Failed to RLP encode state")
	}

	err = stub.PutState("orders", newState)
	if err != nil {
		return shim.Error("Failed to update state of simpleshares")
	}

	// Notify listeners that an event "orderEvent" have been executed
	err = stub.SetEvent("orderEvent", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func main() {
	// Start the chaincode and make it ready for futures requests
	err := shim.Start(new(SimpleShares))
	if err != nil {
		fmt.Printf("Error starting SimpleShare platform: %s", err)
	}
}
