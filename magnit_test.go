package main

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func Test_Init(t *testing.T) {
	magnitCC := new(MAGNIT_CC)
	mockStub := shim.NewMockStub("mockstub", magnitCC)
	txId := "mockTxID"

	mockStub.MockTransactionStart(txId)
	response := magnitCC.Init(mockStub)
	mockStub.MockTransactionEnd(txId)
	if s := response.GetStatus(); s != 200 {
		fmt.Println("Init test failed")
		t.FailNow()
	}
}

/*
* TestInvokeInitTrashTrace simulates an initTrashTrace transaction on the TrashTrace cahincode
 */
func TestInvokeInitMagnit(t *testing.T) {
	fmt.Println("Entering TestInvokeInitMagnit")

	// Instantiate mockStub using TrashTraceDemo as the target chaincode to unit test
	stub := shim.NewMockStub("mockStub", new(MAGNIT_CC))
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}

	var modelName = "Model1234"
	var uploadOgr = "Org1234"
	// Here we perform a "mock invoke" to invoke the function "initmodel" method with associated parameters
	// The first parameter is the function we are invoking

	// data model for initial state - create container
	//      0
	//   "cnt1234"
	result := stub.MockInvoke("001",
		[][]byte{[]byte("initmodel"),
			[]byte(modelName), []byte(uploadOgr)})
	fmt.Println("Status: " + fmt.Sprint(result.GetStatus()))
	// We expect a shim.ok if all goes well
	if result.Status != shim.OK {
		t.Fatalf("Expected unauthorized user error to be returned")
	}

	// here we validate we can retrieve the object we just committed by modelID
	valAsbytes, err := stub.GetState("Model1")
	if err != nil {
		t.Errorf("Failed to get state for Container " + "Model1")
	} else if valAsbytes == nil {
		t.Errorf("Container does not exist:" + "Model1")
	}

}
