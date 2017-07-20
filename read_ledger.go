// =================================================
// AssetChain v0.1 - read object and history
// =================================================

package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ============================================================================================================================
// Read - read a generic variable from ledger
//
// Shows Off GetState() - reading a key/value from the ledger
//
// Inputs - Array of strings
//  0
//  key
//  "abc"
// 
// Returns - string
// ============================================================================================================================
func read(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, jsonResp string
	var err error
	fmt.Println("starting read")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting key of the var to query")
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)           //get the var from ledger
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(jsonResp)
	}

	fmt.Println("- end read")
	return shim.Success(valAsbytes)                  //send it onward
}

// ============================================================================================================================
// Get everything we need (tickets + employee + asset)
//
// Inputs - none
//
// Returns: json array with tickets, employees and assets
// ============================================================================================================================
func read_everything(stub shim.ChaincodeStubInterface) pb.Response {
	type Everything struct {
		Tickets   []Ticket   `json:"tickets"`
		Employees  []Employee  `json:"employee"`
		Assets  []IBM_Asset  `json:"ibmasset"`
	}
	var everything Everything

	// ---- Get All Tickets ---- //
	resultsIterator, err := stub.GetStateByRange("0", "9999999999999999999")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()
	
	for resultsIterator.HasNext() {
		pointer, err := resultsIterator.Next()
		queryKeyAsStr, queryValAsBytes := pointer.GetKey(), pointer.GetValue()
		fmt.Println("queryKeyAsStr =" + queryKeyAsStr)
		if err != nil {
			return shim.Error(err.Error())
		}

		fmt.Println("on ticket id - ", queryKeyAsStr)
		var ticket Ticket
		json.Unmarshal(queryValAsBytes, &ticket)                  //un stringify it aka JSON.parse()
		everything.Tickets = append(everything.Tickets, ticket)   //add this Ticket to the list
	}
	fmt.Println("Tickets array - ", everything.Tickets)

	// ---- Get All Employees ---- //
	employeeIterator, err := stub.GetStateByRange("0", "9999999999999999999")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer employeeIterator.Close()

	for employeeIterator.HasNext() {
		 pointer, err := employeeIterator.Next()
		 queryKeyAsStr, queryValAsBytes := pointer.GetKey(), pointer.GetValue()
		 fmt.Println("queryKeyAsStr =" + queryKeyAsStr)
		if err != nil {
			return shim.Error(err.Error())
		}
		
		fmt.Println("on employee id - ", queryKeyAsStr)
		var employee Employee
		json.Unmarshal(queryValAsBytes, &employee)                  //un stringify it aka JSON.parse()
		everything.Employees = append(everything.Employees, employee)     //add this Employee to the list
	}
	fmt.Println("Employees array - ", everything.Employees)

	// ---- Get All Assets ---- //
	assetIterator, err := stub.GetStateByRange("0", "9999999999999999999")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer assetIterator.Close()

	for assetIterator.HasNext() {
		pointer, err := assetIterator.Next()
		queryKeyAsStr, queryValAsBytes := pointer.GetKey(), pointer.GetValue()
		if err != nil {
			return shim.Error(err.Error())
		}
		
		fmt.Println("on asset id - ", queryKeyAsStr)
		var ibmasset IBM_Asset
		json.Unmarshal(queryValAsBytes, &ibmasset)                  //un stringify it aka JSON.parse()
		everything.Assets = append(everything.Assets, ibmasset)     //add this asset to the list
	}
	fmt.Println("Assets array - ", everything.Assets)

	//change to array of bytes
	everythingAsBytes, _ := json.Marshal(everything)             //convert to array of bytes
	return shim.Success(everythingAsBytes)
}

// ============================================================================================================================
// Get history of ticket
//
// Shows Off GetHistoryForKey() - reading complete history of a key/value
//
// Inputs - Array of strings
//  0
//  id
//  "m01490985296352SjAyM"
// ============================================================================================================================
func getHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	type AuditHistory struct {
		TxId    string   `json:"txId"`
		Value   Ticket   `json:"value"`
	}
	var history []AuditHistory;
	var ticket Ticket

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ticketId := args[0]
	fmt.Printf("- start getHistoryForTicket: %s\n", ticketId)

	// Get History
	resultsIterator, err := stub.GetHistoryForKey(ticketId)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		pointer, err := resultsIterator.Next()
		txID, historicValue := pointer.GetTxId(), pointer.GetValue()
		if err != nil {
			return shim.Error(err.Error())
		}

		var tx AuditHistory
		tx.TxId = txID                             //copy transaction id over
		json.Unmarshal(historicValue, &ticket)     //un stringify it aka JSON.parse()
		if historicValue == nil {                  //marble has been deleted
			var emptyTicket Ticket
			tx.Value = emptyTicket                 //copy nil marble
		} else {
			json.Unmarshal(historicValue, &ticket) //un stringify it aka JSON.parse()
			tx.Value = ticket                      //copy marble over
		}
		history = append(history, tx)              //add this tx to the list
	}
	fmt.Printf("- getHistoryForTicket returning:\n%s", history)

	//change to array of bytes
	historyAsBytes, _ := json.Marshal(history)     //convert to array of bytes
	return shim.Success(historyAsBytes)
}

// ============================================================================================================================
// Get history of tickets - performs a range query based on the start and end keys provided.
//
// Shows Off GetStateByRange() - reading a multiple key/values from the ledger
//
// Inputs - Array of strings
//       0     ,    1
//   startKey  ,  endKey
//  "ticket1" , "ticket2"
// ============================================================================================================================
func getTicketsByRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	startKey := args[0]
	endKey := args[1]

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		pointer, err := resultsIterator.Next()
		queryResultKey, queryResultValue := pointer.GetKey(), pointer.GetValue()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResultKey)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResultValue))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getTicketsByRange queryResult:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
