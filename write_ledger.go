// =================================================
// AssetChain v0.1 - write objects to ledger
// =================================================

package main

import (
	"encoding/json"
	"fmt"
	_ "strconv"
	_ "strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ============================================================================================================================
// write() - genric write variable into ledger
// 
// Shows Off PutState() - writting a key/value into the ledger
//
// Inputs - Array of strings
//    0   ,    1
//   key  ,  value
//  "abc" , "test"
// ============================================================================================================================
func write(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, value string
	var err error
	fmt.Println("starting write")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2. key of the variable and value to set")
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	key = args[0]                                   //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value))         //write the variable into the ledger
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end write")
	return shim.Success(nil)
}

// ============================================================================================================================
// delete_ticket() - remove a ticket from state and from ticket index
// 
// Shows Off DelState() - "removing"" a key/value from the ledger
//
// Inputs - Array of strings
//      0      ,         1
//     id      ,  authed_by_company
// "m999999999", "united marbles"
// ============================================================================================================================
func delete_ticket(stub shim.ChaincodeStubInterface, args []string) (pb.Response) {
	fmt.Println("starting delete_ticket")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// input sanitation
	err := sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	id := args[0]
	
	// get the object
	ticket, err := get_ticket(stub, id)
	if err != nil{
		fmt.Println("Failed to find ticket by id " + id)
		return shim.Error(err.Error())
	}

	// remove the ticket
	err = stub.DelState(ticket.Ticket_Id)   	 //remove the key from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	fmt.Println("- end delete_ticket")
	return shim.Success(nil)
}

// ============================================================================================================================
// delete_employee() - remove an employee from state and from employee index
// 
// Shows Off DelState() - "removing"" a key/value from the ledger
//
// Inputs - Array of strings
//      0      ,         1
//     id      ,  authed_by_company
// "m999999999", "united marbles"
// ============================================================================================================================
func delete_employee(stub shim.ChaincodeStubInterface, args []string) (pb.Response) {
	fmt.Println("starting delete_employee")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// input sanitation
	err := sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	id := args[0]
	
	// get the object
	employee, err := get_employee(stub, id)
	if err != nil{
		fmt.Println("Failed to find employee by id " + id)
		return shim.Error(err.Error())
	}

	// remove the employee
	err = stub.DelState(employee.Employee_sn) 	//remove the key from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	fmt.Println("- end delete_employee")
	return shim.Success(nil)
}

// ============================================================================================================================
// delete_ibmasset() - remove an IBM_Asset from state and from IBM_Asset index
// 
// Shows Off DelState() - "removing"" a key/value from the ledger
//
// Inputs - Array of strings
//      0      ,         1
//     id      ,  authed_by_company
// "m999999999", "united marbles"
// ============================================================================================================================
func delete_ibmasset(stub shim.ChaincodeStubInterface, args []string) (pb.Response) {
	fmt.Println("starting delete_ibmasset")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// input sanitation
	err := sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	id := args[0]

	// get the object
	ibmasset, err := get_ibmasset(stub, id)
	if err != nil{
		fmt.Println("Failed to find asset by id " + id)
		return shim.Error(err.Error())
	}

	// remove the Asset
	err = stub.DelState(ibmasset.SerialNumber)   	 //remove the key from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	fmt.Println("- end delete_ibmasset")
	return shim.Success(nil)
}

// ============================================================================================================================
// Init Ticket - create a new ticket, store into chaincode state
//
// Shows off building a key's JSON value manually
//
// Inputs - Array of strings
//      0      ,    1  ,  2  ,      3          ,       4
//     ticket_id      ,  color, size,     owner id    ,  authing company
// "m999999999", "blue", "35", "o9999999999999", "united marbles"
// ============================================================================================================================
func init_ticket(stub shim.ChaincodeStubInterface, args []string) (pb.Response) {
	var err error
	fmt.Println("starting init_ticket")

	if len(args) != 16 {
		return shim.Error("Incorrect number of arguments. Expecting 16")
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	ticket_id := args[0]
	description := args[1]
	date := args[2]
	status := args[3]
	ticketowner := args[4]
	assignee := args[5]
	asset := args[6]
	queue := args[7]
	address := args[8]
	descriptionproduct := args[9]
	prod := args[10]
	diagnostic := args[11]
	hardwarepw := args[12]
	ospw := args[13]
	contactphone := args[14]
	contactemail := args[15]

	//check if owner exists
	employee, err := get_employee(stub, ticketowner)
	if err != nil {
		fmt.Println("Failed to find employee - " + ticketowner)
		return shim.Error(err.Error())
	}

	//check if ticket id already exists
	ticket, err := get_ticket(stub, ticket_id)
	if err == nil {
		fmt.Println("This ticket already exists - " + ticket_id)
		fmt.Println(ticket)
		return shim.Error("This ticket already exists - " + ticket_id)  //all stop a ticket by this id exists
	}

	//build the ticket json string manually
	str := `{
		"docType":"ticket", 
		"ticket_id": "` + ticket.Ticket_Id + `", 
		"description": "` + description + `", 
		"date": "` + date + `", 
		"status": "` + status + `",
		"ticketowner": "` + employee.Employee_sn + `",
		"assignee": "` + assignee + `",
		"asset": "` + asset + `",
		"queue": "` + queue + `",
		"address": "` + address + `",
		"descriptionproduct": "` + descriptionproduct + `",
		"prod": "` + prod + `",
		"diagnostic": "` + diagnostic + `",
		"hardwarepw": "` + hardwarepw + `",
		"ospw": "` + ospw + `",
		"contactphone": "` + contactphone + `",
		"contactemail": "` + contactemail + `"
	}`
	err = stub.PutState(ticket.Ticket_Id, []byte(str)) 	//store ticket with id as key
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end init_ticket")
	return shim.Success(nil)
}

// ============================================================================================================================
// Init Employee - create a new employee and store into chaincode state
//
// Shows off building key's value from GoLang Structure
//
// Inputs - Array of Strings
//           0     ,     1   ,   2
//      owner id   , username, company
// "o9999999999999",     bob", "united marbles"
// ============================================================================================================================
func init_employee(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting init_employee")

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	var employee Employee
	employee.ObjectType = "employee"
	employee.Employee_sn =  args[0]
	employee.Email = args[1]
	employee.Fullname = args[2]
	fmt.Println(employee)

	//check if employee already exists
	_, err = get_employee(stub, employee.Employee_sn)
	if err == nil {
		fmt.Println("This employee already exists - " + employee.Employee_sn)
		return shim.Error("This employee already exists - " + employee.Employee_sn)
	}

	//store employee
	employeeAsBytes, _ := json.Marshal(employee)	//convert to array of bytes
	err = stub.PutState(employee.Employee_sn, employeeAsBytes)	  //store owner by its Id
	if err != nil {
		fmt.Println("Could not store employee")
		return shim.Error(err.Error())
	}

	fmt.Println("- end init_employee marble")
	return shim.Success(nil)
}

// ============================================================================================================================
// Init Asset - create a new IBM_Asset and store into chaincode state
//
// Shows off building key's value from GoLang Structure
//
// Inputs - Array of Strings
//           0     ,     1   ,   2
//      owner id   , username, company
// "o9999999999999",     bob", "united marbles"
// ============================================================================================================================
func init_ibmasset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting init_ibmasset")

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	var ibmasset IBM_Asset
	ibmasset.ObjectType = "ibm_asset"
	ibmasset.SerialNumber =  args[0]
	ibmasset.AssetType = args[1]
	ibmasset.Tickets =  args[2]
	ibmasset.Owner = args[3]
	fmt.Println(ibmasset)

	//check if asset already exists
	_, err = get_ibmasset(stub, ibmasset.SerialNumber)
	if err == nil {
		fmt.Println("This asset already exists - " + ibmasset.SerialNumber)
		return shim.Error("This asset already exists - " + ibmasset.SerialNumber)
	}

	//store asset
	ibmassetAsBytes, _ := json.Marshal(ibmasset)	//convert to array of bytes
	err = stub.PutState(ibmasset.SerialNumber, ibmassetAsBytes)	  //store asset by its serial number
	if err != nil {
		fmt.Println("Could not store asset")
		return shim.Error(err.Error())
	}

	fmt.Println("- end init_employee asset")
	return shim.Success(nil)
}

// ============================================================================================================================
// Set Assignee on Ticket
//
// Shows off GetState() and PutState()
//
// Inputs - Array of Strings
//       0     ,        1      ,        2
//  marble id  ,  to owner id  , company that auth the transfer
// "m999999999", "o99999999999", united_mables" 
// ============================================================================================================================
func set_assignee(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting set_assignee")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	var ticket_id = args[0]
	var assignee = args[1]
	fmt.Println(ticket_id + "->" + assignee)

	// check if user already exists
	employee, err := get_employee(stub, assignee)
	if err != nil {
		return shim.Error("This employee does not exist - " + assignee)
	}

	// get ticket's current state
	ticketAsBytes, err := stub.GetState(ticket_id)
	if err != nil {
		return shim.Error("Failed to get ticket")
	}
	res := Ticket{}
	json.Unmarshal(ticketAsBytes, &res)           //un stringify it aka JSON.parse()

	// set assignee
	res.Assignee.Employee_sn = employee.Employee_sn                   //change the assignee
	jsonAsBytes, _ := json.Marshal(res)           //convert to array of bytes
	err = stub.PutState(args[0], jsonAsBytes)     //rewrite the ticket with id as key
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end set assignee")
	return shim.Success(nil)
}
