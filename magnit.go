package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

//  Chaincode implementation
type magnitCC struct {
}

// Model data struct
type Model struct {
	// {"Exported": ""}
	ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	ModelID    string `json:"ModelID"`
	ModelName  string `json:"ModelName"`
	UploadOrg  string `json:"UploadOrg"`
}

// AgreementCounterNO need to Agreement ID =  Agreement + AgreementCounterNO
type AgreementCounterNO struct {
	Counter int `json:"counter"`
}

// ModelCounterNO need to  Model ID =  Model + ModelNO
type ModelCounterNO struct {
	Counter int `json:"counter"`
}

// CounterNO for counts....
type CounterNO struct {
	Counter int `json:"counter"`
}

// Agreement data struct
type Agreement struct {
	// {"Exported": ""}
	ObjectType               string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	AgreementID              string `json:"AgreementID"`
	AgreementName            string `json:"AgreementName"`
	AgreementModelID         string `json:"AgreementModelID"`            // model id
	AgreementModelCountUse   string `json:"Agreement_model_account_use"` // model Agreement_model_account_use
	AgreementModelCurrentUse string `json:"AgreementModelCurrentUse"`    // model AgreementModelCurrentUse
	AgreementIssuer          string `json:"AgreementIssuer"`             // org name issuer
	AgreementParticipant     string `json:"AgreementParticipant"`        // org name participant
	AgreementCreateTime      string `json:"AgreementCreateTime"`
	AgreementUpdateTime      string `json:"AgreementUpdateTime"`
	AgreementRemark          string `json:"AgreementRemark"`
	AgreementURLImage        string `json:"AgreementURLImage"`
	AgreementStatus          string `json:"AgreementStatus"`
	AgreementHash            string `json:"AgreementHash"`
}

// ===================================================================================
// Main
// ===================================================================================

func main() {
	err := shim.Start(new(magnitCC))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}

// Init Function Executes only on initializing or on updating the chain code
func (t *magnitCC) Init(APIstub shim.ChaincodeStubInterface) peer.Response {

	// Initializing AgreementCounterNOAssetAsBytes Number

	AgreementCounterNOAssetAsBytes, _ := APIstub.GetState("AgreementCounterNO")

	if AgreementCounterNOAssetAsBytes == nil {
		var AgreementCounter = AgreementCounterNO{Counter: 0}

		AgreementCounterAsBytes, _ := json.Marshal(AgreementCounter)

		err := APIstub.PutState("AgreementCounterNO", AgreementCounterAsBytes)

		if err != nil {

			return shim.Error(fmt.Sprintf("Failed to Intitate AgreementCounterNO"))

		}

	} else {
		AgreementCounter := getCounter(APIstub, "AgreementCounterNO")
		fmt.Printf("AgreementCounterNO is %d", AgreementCounter)

	}

	// Initializing ModelCounterNOAssetAsBytes Number

	ModelCounterNOAssetAsBytes, _ := APIstub.GetState("ModelCounterNO")

	if ModelCounterNOAssetAsBytes == nil {
		var ModelCounter = ModelCounterNO{Counter: 0}

		ModelCounterAsBytes, _ := json.Marshal(ModelCounter)

		err := APIstub.PutState("ModelCounterNO", ModelCounterAsBytes)

		if err != nil {

			return shim.Error(fmt.Sprintf("Failed to Intitate ModelCounterNO"))

		}

	} else {
		ModelCounter := getCounter(APIstub, "ModelCounterNO")
		fmt.Printf("ModelCounterNO is %d", ModelCounter)
	}

	return shim.Success(nil)
}

//getCounter to the latest value of the counter based on the Asset Type provided as input parameter
func getCounter(APIstub shim.ChaincodeStubInterface, AssetType string) int {
	counterAsBytes, _ := APIstub.GetState(AssetType)
	counterAsset := CounterNO{}

	json.Unmarshal(counterAsBytes, &counterAsset)
	fmt.Println("Counter Current Value: ", counterAsset.Counter)
	fmt.Println("for Asset Type:  ", AssetType)

	return counterAsset.Counter
}

//incrementCounter to the increase value of the counter based on the Asset Type provided as input parameter by 1
func incrementCounter(APIstub shim.ChaincodeStubInterface, AssetType string) int {
	counterAsBytes, _ := APIstub.GetState(AssetType)
	counterAsset := CounterNO{}

	json.Unmarshal(counterAsBytes, &counterAsset)
	counterAsset.Counter++
	counterAsBytes, _ = json.Marshal(counterAsset)

	err := APIstub.PutState(AssetType, counterAsBytes)
	if err != nil {

		fmt.Println("Failed to Increment Counter")

	}
	return counterAsset.Counter
}

//updateCounter to the increase value of the counter based on the Asset Type provided as input parameter by NewCount value provided as input
func updateCounter(APIstub shim.ChaincodeStubInterface, AssetType string, NewCount int) int {
	counterAsBytes, _ := APIstub.GetState(AssetType)
	counterAsset := CounterNO{}

	json.Unmarshal(counterAsBytes, &counterAsset)
	counterAsset.Counter = NewCount
	counterAsBytes, _ = json.Marshal(counterAsset)

	fmt.Println("in updateCounter for asset Type :", AssetType)
	fmt.Println("new Count is  :", NewCount)

	err := APIstub.PutState(AssetType, counterAsBytes)
	if err != nil {

		fmt.Println("Failed to Increment Counter")

		return -1

	}
	return counterAsset.Counter
}

// GetTxTimestampChannel Function gets the Transaction time when the chain code was executed it remains same on all the peers where chaincode executes
func (t *magnitCC) GetTxTimestampChannel(APIstub shim.ChaincodeStubInterface) (string, error) {
	txTimeAsPtr, err := APIstub.GetTxTimestamp()
	if err != nil {
		fmt.Printf("Returning error in TimeStamp \n")
		return "Error", err
	}
	fmt.Printf("\t returned value from APIstub: %v\n", txTimeAsPtr)
	timeStr := time.Unix(txTimeAsPtr.Seconds, int64(txTimeAsPtr.Nanos)).String()

	return timeStr, nil
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *magnitCC) Invoke(APIstub shim.ChaincodeStubInterface) peer.Response {
	function, args := APIstub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "initmodel" { //create a new model or student
		return t.initmodel(APIstub, args)
	} else if function == "queryByModelID" { // query a model by id, stupid name - -!
		return t.queryByModelID(APIstub, args)
	} else if function == "insertAgreementinfo" { //insert a Agreement
		return t.insertAgreementinfo(APIstub, args)
	} else if function == "queryByAgreementID" { // query a Agreement
		return t.queryByAgreementID(APIstub, args)
	} else if function == "queryModelByAgreementID" { // delete Model or Agreement
		return t.queryModelByAgreementID(APIstub, args)
	} else if function == "getHistoryForRecord" { //query hisitory of one key for the record
		return t.getHistoryForRecord(APIstub, args)
	} else if function == "queryAllAgreements" { // query all of all agrtements
		return t.queryAllAgreements(APIstub, args)
	} else if function == "queryAllAsset" { // query all of all agrtements
		return t.queryAllAsset(APIstub, args)
	} else if function == "approveAgreement" { // change status to approved
		return t.approveAgreement(APIstub, args)
	} else if function == "del" { // delete Model or Agreement
		return t.del(APIstub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ===========================================================
// del
// ===========================================================
func (t *magnitCC) del(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	id := args[0]
	err := APIstub.DelState(id)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// ============================================================
// initmodel - create a new model, store into chaincode state
// ============================================================
func (t *magnitCC) initmodel(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	if len(args[0]) <= 0 {
		return shim.Error("Model Name argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("Name of organization wich uploaded the model must be a non-empty string")
	}

	// ==== Input sanitation ====
	ModelName := args[0]
	UploadOrg := args[1]

	fmt.Println("- start init model for name ", ModelName)

	ModelCounterNO := getCounter(APIstub, "ModelCounterNO")
	ModelCounterNO++

	ModelID := "Model" + strconv.Itoa(ModelCounterNO)

	// ==== Check if model already exists ====
	modelAsBytes, err := APIstub.GetPrivateData("collectionModel", ModelID)
	if err != nil {
		return shim.Error("Failed to get model: " + err.Error())
	} else if modelAsBytes != nil {
		fmt.Println("This model already exists: " + ModelID)
		return shim.Error("This model already exists: " + ModelID)
	}

	// ==== Create model object and marshal to JSON ====
	objectType := "model"
	Model := &Model{objectType, ModelID, ModelName, UploadOrg}
	ModelJSONasBytes, err := json.Marshal(Model)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save model to state ===
	err = APIstub.PutPrivateData("collectionModel", ModelID, ModelJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	incCount := incrementCounter(APIstub, "ModelCounterNO")

	// ==== model saved and indexed. Return success ====
	fmt.Println("- end init, new count for Model is  :\n", incCount)
	return shim.Success(nil)

}

// ===============================================
// queryByAgreementID - read a Agreement from chaincode state
// ===============================================
func (t *magnitCC) queryByAgreementID(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	var AgreementID, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting AgreementID of the Agreement to query")
	}

	AgreementID = args[0]
	valAsbytes, err := APIstub.GetState(AgreementID) //get the Agreement from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + AgreementID + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Agreement does not exist: " + AgreementID + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// =================================================================
// queryModelByAgreementID - read a Agreement from chaincode state
// =================================================================
func (t *magnitCC) queryModelByAgreementID(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	var AgreementID, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting AgreementID of the Agreement to query")
	}

	AgreementID = args[0]

	valAsbytes, err := APIstub.GetState(AgreementID) //get the Agreement from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + AgreementID + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Agreement does not exist: " + AgreementID + "\"}"
		return shim.Error(jsonResp)
	}

	Agreement := &Agreement{}
	err = json.Unmarshal([]byte(valAsbytes), &Agreement)
	if err != nil {
		return shim.Error(err.Error())
	}
	// check of uses
	countUse, err := strconv.Atoi(Agreement.AgreementModelCountUse)
	if err != nil {
		fmt.Printf("Can't convert to int AgreementAsset.AgreementModelCountUse %s\n", Agreement.AgreementModelCountUse)
	}
	currentCount, err := strconv.Atoi(Agreement.AgreementModelCurrentUse)
	if err != nil {
		fmt.Printf("Can't convert to int AgreementAsset.AgreementModelCurrentUse %s\n", Agreement.AgreementModelCurrentUse)
	}

	// return error if  currentCount > countUse

	if currentCount >= countUse {
		jsonResp = "{\"Error\":\"Failed to get the Model - model_current_count_query: " + " Вы достигли лимита разрешенных запросов: " + Agreement.AgreementModelCountUse + "\"}"
		return shim.Error(jsonResp)
	}

	Agreement.AgreementModelCurrentUse = strconv.Itoa(currentCount)

	eventPayload := "Agreement with ID " + Agreement.AgreementID + " was selected"
	payloadAsBytes := []byte(eventPayload)
	eventErr := APIstub.SetEvent("queryEvent", payloadAsBytes)
	if eventErr != nil {
		return shim.Error(fmt.Sprintf("Failed to emit event"))
	}
	fmt.Println("Event: Agrrement with ID " + Agreement.AgreementID + " was selected")
	// increase count of uses of Agreement
	output := t.updateAgreement(APIstub, *Agreement)

	if output != "Success" {
		return shim.Error(err.Error())
	}

	modelAsbytes, err := APIstub.GetPrivateData("collectionModel", Agreement.AgreementModelID) //get the model from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + Agreement.AgreementModelID + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"model does not exist: " + Agreement.AgreementModelID + "\"}"
		return shim.Error(jsonResp)
	}

	Model := &Model{}
	err = json.Unmarshal([]byte(modelAsbytes), &Model)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(modelAsbytes)
}

// ===============================================================
// insertAgreementinfo - insert A new Agreement information
//
//AgreementID
//AgreementName
//AgreementModelID
//AgreementModelCountUse
//AgreementModelCurrentUse
//AgreementIssuer
//AgreementParticipant
//AgreementCreateTime
//AgreementUpdateTime
//AgreementRemark
//AgreementURLImage
//AgreementStatus string
//AgreementHash string
// ===============================================================
func (t *magnitCC) insertAgreementinfo(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 9 {
		return shim.Error("##Incorrect number of arguments. expecting 9 args")
	}

	AgreementName := args[0]
	AgreementModelID := args[1]
	AgreementModelCountUse := args[2]
	AgreementModelCurrentUse := "0"
	AgreementIssuer := args[3]
	AgreementParticipant := args[4]
	AgreementRemark := args[5]
	AgreementURLImage := args[6]
	AgreementStatus := args[7]
	AgreementHash := args[8]

	AgreementCounterNO := getCounter(APIstub, "AgreementCounterNO")
	AgreementCounterNO++

	AgreementID := "Agreement" + strconv.Itoa(AgreementCounterNO)

	fmt.Println("###start insertAgreementinfo ID: ", AgreementID)

	AgreementCreateTime, errTx := t.GetTxTimestampChannel(APIstub)
	if errTx != nil {
		return shim.Error("Returning error")
	}

	AgreementUpdateTime, errTx := t.GetTxTimestampChannel(APIstub)
	if errTx != nil {
		return shim.Error("Returning error")
	}

	// check if model exists
	valAsBytes, err := APIstub.GetState(AgreementModelID)
	if err != nil {
		return shim.Error("Failed to get model:" + AgreementModelID + "," + err.Error())
	} else if valAsBytes == nil {
		fmt.Println("Model id does not exist:[" + AgreementModelID + "]")
		return shim.Error("Model id does not exist" + AgreementModelID)
	}

	//check if Agreement exist
	//AgreementAsBytes, err := APIstub.GetState(AgreementID)
	//if err != nil {
	//	return shim.Error("Failed to get Agreement:" + AgreementID + "," + err.Error())
	//}
	//if AgreementAsBytes == nil {
	//	fmt.Println("starting to insert Agreement")
	//}

	objectType := "Agreement"
	Agreement := &Agreement{objectType, AgreementID, AgreementName, AgreementModelID, AgreementModelCountUse, AgreementModelCurrentUse, AgreementIssuer, AgreementParticipant, AgreementCreateTime, AgreementUpdateTime, AgreementRemark, AgreementURLImage, AgreementStatus, AgreementHash}
	AgreementJSONasBytes, err := json.Marshal(Agreement)
	if err != nil {
		return shim.Error(err.Error())
	}

	//AgreementJSONasBytes, _ := json.Marshal(batchToUpdate)
	err = APIstub.PutState(AgreementID, AgreementJSONasBytes) //insert the Agreement
	if err != nil {
		return shim.Error(err.Error())
	}

	incCount := incrementCounter(APIstub, "AgreementCounterNO")

	// ==== modelagreement saved and indexed. Return success ====

	eventPayload := "Agreement with ID " + AgreementID + " was issued and ready to confirm"
	payloadAsBytes := []byte(eventPayload)
	eventErr := APIstub.SetEvent("newAgreementEvent", payloadAsBytes)
	if eventErr != nil {
		return shim.Error(fmt.Sprintf("Failed to emit event"))
	}
	fmt.Println("Event: Agrrement with ID " + Agreement.AgreementID + " was selected")

	fmt.Println("------  end insertAgreementinfo  (success) incCount: ", incCount)
	fmt.Println("------  end insertAgreementinfo  (success) AgreementID: " + AgreementID)
	return shim.Success(nil)
}

// ===============================================================
// queryByModelId - read data for one model from chaincode state
// ===============================================================
func (t *magnitCC) queryByModelID(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	var ModelID, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting ModelID to query")
	}

	ModelID = args[0]
	valAsbytes, err := APIstub.GetPrivateData("collectionModel", ModelID) //get the model from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + ModelID + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"model does not exist: " + ModelID + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)

}

// ==========================================================
// approveAgreement - update status of Agreement to agrg[1]
// ==========================================================
func (t *magnitCC) approveAgreement(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	var err error
	// check args
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	if len(args[0]) <= 0 {
		return shim.Error("AgreementID must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("Status of agreement must be a non-empty string")
	}

	AgreementID := args[0]
	status := args[1]

	updateTime, errTx := t.GetTxTimestampChannel(APIstub)
	if errTx != nil {
		return shim.Error("Returning error")
	}

	valAsbytes, err := APIstub.GetState(AgreementID) //get the Agreement from chaincode state
	if err != nil {
		return shim.Error(err.Error())
	} else if valAsbytes == nil {
		return shim.Error("Agreement not exist")
	}

	Agreement := &Agreement{}
	err = json.Unmarshal([]byte(valAsbytes), &Agreement)
	if err != nil {
		return shim.Error(err.Error())
	}
	Agreement.AgreementStatus = status
	Agreement.AgreementUpdateTime = updateTime

	valAsbytes, err = json.Marshal(Agreement)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = APIstub.PutState(AgreementID, valAsbytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("update success")
	return shim.Success(nil)

}

// ==============================================================
// updateAgreement update count of use in the Agreement count++
// ==============================================================
func (t *magnitCC) updateAgreement(APIstub shim.ChaincodeStubInterface, AgreementAsset Agreement) string {

	fmt.Printf("In update AgreementId %s -- array %v", AgreementAsset.AgreementID, AgreementAsset)

	updateTime, errTx := t.GetTxTimestampChannel(APIstub)
	if errTx != nil {
		return "GetTxTimestampChannel returning error"
	}

	countUse, err := strconv.Atoi(AgreementAsset.AgreementModelCurrentUse)
	if err != nil {
		fmt.Printf("Can't convert to int AgreementAsset.AgreementModelCountUse %s\n", AgreementAsset.AgreementModelCountUse)
	}

	countUse++

	AgreementAsset.AgreementModelCurrentUse = strconv.Itoa(countUse)
	AgreementAsset.AgreementUpdateTime = updateTime

	fmt.Printf("Increase count:%s for %s", AgreementAsset.AgreementModelCurrentUse, AgreementAsset.AgreementID)

	AgreementUp := &Agreement{AgreementAsset.ObjectType, AgreementAsset.AgreementID, AgreementAsset.AgreementName, AgreementAsset.AgreementModelID, AgreementAsset.AgreementModelCountUse, AgreementAsset.AgreementModelCurrentUse, AgreementAsset.AgreementIssuer, AgreementAsset.AgreementParticipant, AgreementAsset.AgreementCreateTime, AgreementAsset.AgreementUpdateTime, AgreementAsset.AgreementRemark, AgreementAsset.AgreementURLImage, AgreementAsset.AgreementStatus, AgreementAsset.AgreementHash}

	valJSONasBytes, err := json.Marshal(AgreementUp)
	if err != nil {
		return "error: Failed to json.Marshal valJSONasBytes"
	}

	err = APIstub.PutState(AgreementAsset.AgreementID, valJSONasBytes)
	if err != nil {
		return "error: Failed to PutState AgreementID, AgreementAsset"
	}

	fmt.Println("Succes update Agreement: ", valJSONasBytes)
	return "Success"
}

// ===============================================
// queryAllAgreements in the channel
// ===============================================

func (t *magnitCC) queryAllAgreements(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	queryString := "{\"selector\":{\"docType\":\"Agreement\"}}"

	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	//return shim.Success()
	return shim.Success(queryResults)
}

// query all assets
func (t *magnitCC) queryAllAsset(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	startKey := ""

	endKey := ""

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)

	if err != nil {

		return shim.Error(err.Error())

	}

	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults

	var buffer bytes.Buffer

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false

	for resultsIterator.HasNext() {

		queryResponse, err := resultsIterator.Next()
		// respValue := string(queryResponse.Value)
		if err != nil {

			return shim.Error(err.Error())

		}

		// Add a comma before array members, suppress it for the first array member

		if bArrayMemberAlreadyWritten == true {

			buffer.WriteString(",")

		}

		buffer.WriteString("{\"Key\":")

		buffer.WriteString("\"")

		buffer.WriteString(queryResponse.Key)

		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")

		// Record is a JSON object, so we write as-is

		buffer.WriteString(string(queryResponse.Value))

		buffer.WriteString("}")

		bArrayMemberAlreadyWritten = true

	}

	buffer.WriteString("]")

	fmt.Printf("- queryAllAssets:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())

}

// ===============================================
// createIndexformodel - build an index for an model to
// ===============================================
func (t *magnitCC) createIndexformodel(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	//return shim.Success()
	return shim.Success(nil)
}

// ===========================================================================================
// getHistoryForRecord returns the historical state transitions for a given key of a record
// ===========================================================================================
func (t *magnitCC) getHistoryForRecord(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	recordKey := args[0]

	fmt.Printf("- start getHistoryForRecord: %s\n", recordKey)

	resultsIterator, err := APIstub.GetHistoryForKey(recordKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the key/value pair
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON goods)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForRecord returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// ========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(APIstub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := APIstub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}
