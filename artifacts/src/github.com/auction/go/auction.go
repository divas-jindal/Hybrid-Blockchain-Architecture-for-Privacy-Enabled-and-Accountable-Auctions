package main

import (
	// "bytes"
	"encoding/json"
	"fmt"
	"strconv"
	// "time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/flogging"

	// "github.com/hyperledger/fabric-chaincode-go/pkg/cid"
)

// SmartContract Define the Smart Contract structure
var err error
// SmartContract Define the Smart Contract structure
type SmartContract struct {
}

// to intialize structs here
type State struct {
	BidStarted string `json:"bidStarted"`
	RevealStarted string `json:"revealStarted"`
}

type Commitment struct {
	UserID int `json:"userID"`
	RandomValue int `json:"randomValue"`
	BidValue int `json:"bidValue"`
}

type PublicCommitment struct {
	UserID int `json:"userID"`
	RandomValue int `json:"randomValue"`
}

type Bid struct {
	HashedBidValue int `json:"hashedBidValue"`
}

type WinnerCandidates struct {
	HighestBidder string `json:"highestBidder"`
	HighestBid int `json:"highestBid"`
	SecondHighestBidder string `json:"secondHighestBidder"`
	SecondHighestBid int `json:"secondHighestBid"`
}

type Winner struct {
	Winner string `json:"winner"`
	Amount int `json:"amount"`
}

// Init ;  Method for initializing smart contract
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	state := State{BidStarted: "false", RevealStarted: "false"}
	stateAsBytes, _ := json.Marshal(state)
	APIstub.PutState("State", stateAsBytes)


	commitment := Commitment{RandomValue: 0, BidValue: 0}
	commitmentAsBytes, _ := json.Marshal(commitment)
	APIstub.PutState("Commitment", commitmentAsBytes)
	return shim.Success(nil)
}

var logger = flogging.MustGetLogger("auction_cc")

// Invoke :  Method for INVOKING smart contract
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()
	logger.Infof("Function name is:  %d", function)
	logger.Infof("Args length is : %d", len(args))

	// to write code for each function
	switch function {
	case "bidstart":
		return s.BidStart(APIstub)
	case "revealstart":
		return s.RevealStart(APIstub)
	case "makecommitmentorg":
		return s.MakeCommitmentOrg(APIstub, args)
	case "revealcommitment":
		return s.RevealCommitment(APIstub, args)
	default:
		return shim.Error("Invalid Smart Contract function name.")
	}

}

func (s *SmartContract) BidStart(APIstub shim.ChaincodeStubInterface) sc.Response{
	// check role for auctioneer

	stateAsBytes, _ := APIstub.GetState("State")
	state := State{}

	json.Unmarshal(stateAsBytes, &state)
	state.BidStarted = "true"
	state.RevealStarted = "false"

	stateAsBytes, _ = json.Marshal(state)
	APIstub.PutState("State", stateAsBytes)
	return shim.Success(stateAsBytes)
	// to send to all members of channel to make commitment
	// timebound
}

func (s *SmartContract) RevealStart(APIstub shim.ChaincodeStubInterface) sc.Response{
	// check role for auctioneer

	stateAsBytes, _ := APIstub.GetState("State")
	state := State{}

	json.Unmarshal(stateAsBytes, &state)
	state.BidStarted = "false"
	state.RevealStarted = "true"

	stateAsBytes, _ = json.Marshal(state)
	APIstub.PutState("State", stateAsBytes)
	return shim.Success(stateAsBytes)
	// to send to all members of channel to make commitment
	// timebound
}

func (s* SmartContract)  MakeCommitmentOrg(APIstub shim.ChaincodeStubInterface, args []string) sc.Response{

	// check role for participant

	commitment := Commitment{}

	// id, err := cid.GetID(APIstub)
	commitment.UserID, err =  strconv.Atoi(args[0])
	commitment.RandomValue, err = strconv.Atoi(args[1])
	commitment.BidValue, err = strconv.Atoi(args[2])

	commitmentAsBytes,_ := json.Marshal(commitment)

	APIstub.PutState(args[0], commitmentAsBytes)
	// APIstub.PutState(id, commitmentAsBytes)

	return shim.Success(commitmentAsBytes)

}

func (s* SmartContract) RevealCommitment(APIstub shim.ChaincodeStubInterface, args []string) sc.Response{
	// check role for auctioneer
	commitmentAsBytes, _ := APIstub.GetState(args[0])
	commitment := Commitment{}
	json.Unmarshal(commitmentAsBytes, &commitment)
	var randomValue int
	randomValue, err = strconv.Atoi(args[1])
	if(commitment.RandomValue == randomValue) { // && OneWayHash(commitment.BidValue) == randomValue) {
		// CheckIfWinner(APIstub, args[0], commitment.BidValue);
		return shim.Success(commitmentAsBytes);
	} else {
		return shim.Error("Fraud")	
	}
}

// func CheckIfWinner(APIstub shim.ChaincodeStubInterface, UserID string, BidValue int) {
// 		winnerAsBytes, _ := APIstub.GetPrivateData("auctionPrivateDetails", "WinnerDetails")
// 		winner := Winner{}

// 		json.Unmarshal(winnerAsBytes, &winner)
// 		var bidValue int
// 		bidValue, err = strconv.Atoi(args[2])
// 		if (bidValue > winner.HighestBid) {
// 			winner.SecondHighestBidder = winner.HighestBidder
// 			winner.SecondHighestBid = winner.HighestBid
// 			winner.HighestBid = bidValue;
// 			winner.HighestBidder = args[0]
// 		} else if (bidValue > winner.SecondHighestBid) {
// 			winner.SecondHighestBid = bidValue
// 			winner.SecondHighestBidder = args[0]
// 		}

// 		winnerAsBytes, _  = json.Marshal(winner)
// 		APIstub.PutPrivateData("auctionPrivateData", "WinnerDetails", winnerAsBytes)
// }

func main() {

	// Create a new Smart Contract
	err = shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}


