package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	// "time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/flogging"

	// "github.com/hyperledger/fabric-chaincode-go/pkg/cid"
)

var err error
// SmartContract Define the Smart Contract structure
type SmartContract struct {
}

// to intialize structs here
type State struct {
	BidStarted string `json:"bidStarted"`
	RevealStarted string `json:"revealStarted"`
	Breach string `json:"breach"`
}

type Commitment struct {
	UserID int `json:"userID"`
	RandomValue int `json:"randomValue"`
	BidValue int `json:"bidValue"`
}

type PublicCommitment struct {
	UserID int `json:"userID"`
	RandomValue int `json:"randomValue"`
	BidValue int `json:"bidValue"`
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

	commitment := Commitment{}
	commitmentAsBytes,_ := json.Marshal(commitment)
	APIstub.PutPrivateData("VotingPrivateData", "myCommitment", commitmentAsBytes)
	return shim.Success(nil)
}

var logger = flogging.MustGetLogger("voting_cc")

// Invoke :  Method for INVOKING smart contract
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()
	logger.Infof("Function name is:  %d", function)
	logger.Infof("Args length is : %d", len(args))

	// to write code for each function
	switch function {
	case "bidfromOrg2":
		return s.BidFromOrg2(APIstub, args)
	case "bidfromOrg3":
		return s.BidFromOrg3(APIstub, args)
	case "bidrecvOrg":
		return s.BidRecvOrg(APIstub, args)
	case "declareWinner":
		return s.DeclareWinner(APIstub, args)
	case "checkIfWinner":
		return s.CheckIfWinner(APIstub, args)
	case "breach":
		return s.Breach(APIstub, args)
	case "bringToPublic":
		return s.BringToPublic(APIstub, args)
	case "checkAll":
		return s.CheckAll(APIstub, args)
	default:
		return shim.Error("Invalid Smart Contract function name.")
	}

}

func (s *SmartContract) BidFromOrg2(APIstub shim.ChaincodeStubInterface, args []string) sc.Response{
	publicCommitment := PublicCommitment{}

	commitmentAsBytes,_ := APIstub.GetPrivateData("votingPrivateDetailsOrg2", "myCommitment")
	commitment := Commitment{}
	if(commitmentAsBytes != nil) {
		json.Unmarshal(commitmentAsBytes, &commitment)
	}

	// id, err := cid.GetID(APIstub)
	commitment.UserID, err =  strconv.Atoi(args[0])
	commitment.RandomValue, err = strconv.Atoi(args[1])
	commitment.BidValue, err = strconv.Atoi(args[2])

	publicCommitment.UserID, err = strconv.Atoi(args[0])
	publicCommitment.RandomValue, err = strconv.Atoi(args[1])

	commitmentAsBytes,_ = json.Marshal(commitment)
	publicCommitmentAsBytes, _ := json.Marshal(publicCommitment)

	APIstub.PutPrivateData("votingPrivateDetailsOrg2", "myCommitment", commitmentAsBytes)
	APIstub.PutState(args[0], publicCommitmentAsBytes)
	// APIstub.PutState(id, publicCommitmentAsBytes)

	return shim.Success(nil)
	
}

func (s *SmartContract) BidFromOrg3(APIstub shim.ChaincodeStubInterface, args []string) sc.Response{
	publicCommitment := PublicCommitment{}

	commitmentAsBytes,_ := APIstub.GetPrivateData("votingPrivateDetailsOrg3", "myCommitment")
	commitment := Commitment{}
	json.Unmarshal(commitmentAsBytes, &commitment)

	// id, err := cid.GetID(APIstub)
	commitment.UserID, err =  strconv.Atoi(args[0])
	commitment.RandomValue, err = strconv.Atoi(args[1])
	commitment.BidValue, err = strconv.Atoi(args[2])

	publicCommitment.UserID, err = strconv.Atoi(args[0])
	publicCommitment.RandomValue, err = strconv.Atoi(args[1])

	commitmentAsBytes,_ = json.Marshal(commitment)
	publicCommitmentAsBytes, _ := json.Marshal(publicCommitment)

	APIstub.PutPrivateData("votingPrivateDetailsOrg3", "myCommitment", commitmentAsBytes)
	APIstub.PutState(args[0], publicCommitmentAsBytes)
	// APIstub.PutState(id, publicCommitmentAsBytes)

	return shim.Success(nil)
	
}

func (s *SmartContract) BidRecvOrg(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	publicCommitmentAsBytes,_ := APIstub.GetState(args[0])
	return shim.Success(publicCommitmentAsBytes)
	
}

func (s *SmartContract) CheckIfWinner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response{
	winnerCandidatesAsBytes, _ := APIstub.GetPrivateData("auctionPrivateDetails", "WinnerCandidatesDetails")
	winnerCandidates := WinnerCandidates{}

	if(winnerCandidatesAsBytes != nil) {
		json.Unmarshal(winnerCandidatesAsBytes, &winnerCandidates)
	} 
	var bidValue int
	bidValue, err = strconv.Atoi(args[2])
	if (len(winnerCandidates.HighestBidder) == 0 || bidValue > winnerCandidates.HighestBid) {
		winnerCandidates.SecondHighestBidder = winnerCandidates.HighestBidder
		winnerCandidates.SecondHighestBid = winnerCandidates.HighestBid
		winnerCandidates.HighestBid = bidValue
		winnerCandidates.HighestBidder = args[0]
		winnerCandidatesAsBytes, _  = json.Marshal(winnerCandidates)
		APIstub.PutPrivateData("auctionPrivateData", "WinnerCandidatesDetails", winnerCandidatesAsBytes)
		var buffer bytes.Buffer
		buffer.WriteString("[{\"Msg\":\"Bid Increased Event\"}]")
		return shim.Success(buffer.Bytes())
	} else if (len(winnerCandidates.SecondHighestBidder) == 0 || bidValue > winnerCandidates.SecondHighestBid) {
		winnerCandidates.SecondHighestBid = bidValue
		winnerCandidates.SecondHighestBidder = args[0]
	}

	winnerCandidatesAsBytes, _  = json.Marshal(winnerCandidates)
	APIstub.PutPrivateData("auctionPrivateData", "WinnerCandidatesDetails", winnerCandidatesAsBytes)
	return shim.Success(nil)
}

func (s* SmartContract) DeclareWinner (APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	winnerCandidatesAsBytes, _ := APIstub.GetPrivateData("auctionPrivateDetails", "WinnerCandidatesDetails")
	winnerCandidates := WinnerCandidates{}
	json.Unmarshal(winnerCandidatesAsBytes, &winnerCandidates)

	winner := Winner{}
	winner.Amount = winnerCandidates.SecondHighestBid
	winner.Winner = winnerCandidates.HighestBidder

	winnerAsBytes, _  := json.Marshal(winner)

	APIstub.PutState("Winner", winnerAsBytes)
	return shim.Success(nil)
}


// to be done
func (s*SmartContract) Breach(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	
	//check all commitments
	stateAsBytes, _ := APIstub.GetState("State")
	state := State{}

	json.Unmarshal(stateAsBytes, &state)
	state.Breach = args[0]

	stateAsBytes, _ = json.Marshal(state)
	APIstub.PutState("State", stateAsBytes)
	var buffer bytes.Buffer
	buffer.WriteString("[{\"Msg\":\"Bring commitments to public\"}]")
	return shim.Success(buffer.Bytes())
}

func (s*SmartContract) BringToPublic(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	publicCommitmentAsBytes,_ := APIstub.GetState(args[0])
	publicCommitment := PublicCommitment{}

	json.Unmarshal(publicCommitmentAsBytes, &publicCommitment)
	publicCommitment.BidValue,err = strconv.Atoi(args[1])
	return shim.Success(publicCommitmentAsBytes)
}

func (s*SmartContract) CheckAll(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	startKey := "p111"
	endKey := "p999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	winnerCandidatesAsBytes, _ := APIstub.GetPrivateData("auctionPrivateDetails", "WinnerCandidatesDetails")
	winnerCandidates := WinnerCandidates{}

	if(winnerCandidatesAsBytes != nil) {
		json.Unmarshal(winnerCandidatesAsBytes, &winnerCandidates)
	} 
	var buffer bytes.Buffer

	// bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		publicCommitmentAsBytes := []byte(queryResponse.Value)
		publicCommitment := PublicCommitment{}
		json.Unmarshal(publicCommitmentAsBytes, &publicCommitment)
		if(publicCommitment.BidValue > winnerCandidates.HighestBid) {
			buffer.WriteString("[{\"Msg\":\"Accuse " + winnerCandidates.HighestBidder+"\"}]")
			return shim.Success(buffer.Bytes())
		}
		
	}

	stateAsBytes, _ := APIstub.GetState("State")
	state := State{}

	json.Unmarshal(stateAsBytes, &state)
	buffer.WriteString("[{\"Msg\":\"Penalize " + state.Breach+"\"}]")
	return shim.Success(buffer.Bytes())
}

func main() {

	// Create a new Smart Contract
	err = shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}


