package main

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/clearmatics/autonity/rlp"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// unmarshallState receives state in byte form and returns unmarshalled struct
func unmarshallState(state []byte) (Shares, error) {
	// Unmarshall returned state into struct
	var shares Shares
	err := json.Unmarshal(state, &shares)
	if err != nil {
		return shares, errors.New("Failed unmarshall state")
	}

	return shares, nil
}

// marshallState receives state struct and returns byte format json state
func marshallState(shares Shares) ([]byte, error) {
	newState, err := json.MarshalIndent(shares, "", " ")
	if err != nil {
		return nil, err
	}

	return newState, nil
}

// unmarshallOrderBook receives state in byte form and returns unmarshalled struct
func unmarshallOrderBook(state []byte) (OrderStruct, error) {
	// Unmarshall returned state into struct
	var orders OrderStruct
	err := json.Unmarshal(state, &orders)
	if err != nil {
		return orders, errors.New("Failed unmarshall state")
	}

	return orders, nil
}

// marshallState receives state struct and returns byte format json state
func marshallOrderBook(orders OrderStruct) ([]byte, error) {
	output, err := json.MarshalIndent(orders, "", " ")
	if err != nil {
		return nil, err
	}

	return output, nil
}

// marshallTrades receives state struct and returns byte format json state
func marshallTrades(trades Trades) ([]byte, error) {
	output, err := json.MarshalIndent(trades, "", " ")
	if err != nil {
		return nil, err
	}

	return output, nil
}

//
func newOrder(args []string) (OpenOrder, error) {
	// Create new order from input args
	var order OpenOrder
	order.Organisation = args[1]

	// identify the order type
	if args[2] == "Buy" {
		order.BuySell = "Buy"
		order.Address = common.HexToAddress(args[5])
	} else if args[2] == "Sell" {
		order.BuySell = "Sell"
		order.Address = common.HexToAddress(args[5])
	} else {
		return order, errors.New("Incorrect Buy/Sell Flag")
	}

	// Retrieve amount
	amount, err := strconv.Atoi(args[3])
	if err != nil {
		return order, errors.New("Amount is not integer")
	}
	order.Amount = uint(amount)

	// Retrieve price
	amount, err = strconv.Atoi(args[4])
	if err != nil {
		return order, errors.New("Price is not integer")
	}
	order.Price = uint(amount)

	// Retrieve reference
	order.Ref = args[6]

	return order, nil
}

// matchOrder reviews order book to see if a matching order is included
func matchOrder(orderBook OrderStruct, newOrder OpenOrder) OpenOrder {
	var matchedOrder OpenOrder

	// Match order
	for i := range orderBook.Open {
		// if orderBook.Open[i].Organisation == newOrder.Organisation {
		if checkDetails(orderBook.Open[i], newOrder) {
			matchedOrder = orderBook.Open[i]
			return matchedOrder
		}
	}

	return matchedOrder
}

// orderMatched
func orderMatched(orderBook OrderStruct, matchedOrder OpenOrder, newOrder OpenOrder) OrderStruct {
	// Add matched order to section in db
	var newMatched MatchedOrder
	newMatched.Organisation = matchedOrder.Organisation
	newMatched.Amount = matchedOrder.Amount
	newMatched.Price = matchedOrder.Price
	newMatched.Ref = matchedOrder.Ref

	// Determine seller and buyer
	if newOrder.BuySell == "Buy" {
		newMatched.Recv = newOrder.Address
		newMatched.Send = matchedOrder.Address
	} else {
		newMatched.Send = newOrder.Address
		newMatched.Recv = matchedOrder.Address
	}
	orderBook.Matched = append(orderBook.Matched, newMatched)

	// Remove original order from orderBook
	for i := range orderBook.Open {
		if orderBook.Open[i] == matchedOrder {
			orderBook.Open = append(orderBook.Open[:i], orderBook.Open[i+1:]...)
		}
	}

	return orderBook
}

// checkDetails makes sure that orders have corresponding values
func checkDetails(orderA OpenOrder, orderB OpenOrder) bool {
	var tmpA OpenOrder
	var tmpB OpenOrder

	// Create new temp value
	tmpA = orderA
	tmpB = orderB

	// ensure that orders are not both sell/buy orders
	if tmpA.BuySell == tmpB.BuySell {
		return false
	}

	// Reset values that we do not wish to compare
	tmpA.BuySell = tmpB.BuySell
	tmpA.Address = tmpB.Address
	tmpA.Ref = tmpB.Ref

	if tmpA == tmpB {
		return true
	}

	return false

}

func retrieveOrderBook(stub shim.ChaincodeStubInterface) (OrderStruct, error) {
	var orderBook OrderStruct
	// Retrieve order details
	state, err := stub.GetState("orders")
	if err != nil {
		return orderBook, errors.New("Failed to get state")
	}

	err = rlp.DecodeBytes(state, &orderBook)
	if err != nil {
		return orderBook, errors.New("Failed decode rlp bytes")
	}

	return orderBook, nil
}

func retrieveShareListings(stub shim.ChaincodeStubInterface) (Shares, error) {
	var listings Shares

	// Retrieve share listings
	state, err := stub.GetState("shares")
	if err != nil {
		return listings, errors.New("Failed to get state")
	}

	err = rlp.DecodeBytes(state, &listings)
	if err != nil {
		return listings, errors.New("Failed to get state")
	}

	return listings, nil
}
