package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// Define the EIP712Domain type
type EIP712Domain struct {
	Name              string         `json:"name"`
	Version           string         `json:"version"`
	ChainId           *big.Int       `json:"chainId"`
	VerifyingContract common.Address `json:"verifyingContract"`
}

// Define the Order type
type Order struct {
	Salt         *big.Int       `json:"salt"`
	Maker        common.Address `json:"maker"`
	Receiver     common.Address `json:"receiver"`
	MakerAsset   common.Address `json:"makerAsset"`
	TakerAsset   common.Address `json:"takerAsset"`
	MakingAmount *big.Int       `json:"makingAmount"`
	TakingAmount *big.Int       `json:"takingAmount"`
	MakerTraits  *big.Int       `json:"makerTraits"`
}

func main() {

	// Assuming the `order_data` is provided as in your Python example and converted to Go types
	order_data := map[string]string{
		// keccak256 hash of the extension for salt, guaranteed to include the lower 160 bits of the keccak256 hash of the extension
		"salt":         "38190759697773312974901587002227572529508873579400343413173261273827457248092",
		"makerAsset":   "0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619",
		"takerAsset":   "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359",
		"maker":        "0x2c9b2DBdbA8A9c969Ac24153f5C1c23CB0e63914",
		"receiver":     "0x0000000000000000000000000000000000000000",
		"makingAmount": "200000000000000",
		"takingAmount": "462657",
		// manually edited maker traits, 0 nonce
		"makerTraits": "33471150795161712739625987854073848363835857058898165664301626148231099449344",
	}

	extension := "0x00000203000001220000012200000122000001220000009100000000000000002ad5004c60e16e54d5007c80ce329adde5b51ef5000000000000006859e6260000b401bf92000000000000640ac0866635457d36ab318d0000000000000000000066593d4e7d3a5f55167fd18bd45f0b94f54a968f000000000000000000000000000000000000000000000000000000000000c976bf098c4dba0a061d972ad4499f120902631a95770895ad27ad6b0d952ad5004c60e16e54d5007c80ce329adde5b51ef5000000000000006859e6260000b401bf92000000000000640ac0866635457d36ab318d0000000000000000000066593d4e7d3a5f55167fd18bd45f0b94f54a968f000000000000000000000000000000000000000000000000000000000000c976bf098c4dba0a061d972ad4499f120902631a95770895ad27ad6b0d952ad5004c60e16e54d5007c80ce329adde5b51ef500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000646859e6150ac0866635457d36ab318d000000000000000000000000000066593d4e7d3a5f55167f0000d18bd45f0b94f54a968f0000000000000000000000000000000000000000000000000000000000000000000000000000c976bf098c4dba0a061d0000972ad4499f120902631a000095770895ad27ad6b0d9500000000000000000000000000000000000000000000000000000000000000075dec5a"

	chainIdInt := 137                                     // Polygon (match JS/Python)
	chainId := math.NewHexOrDecimal256(int64(chainIdInt)) // Polygon
	apiKey := os.Getenv("DEV_PORTAL_TOKEN")

	// Set up the domain data
	domainData := apitypes.TypedDataDomain{
		Name:              "1inch Aggregation Router",
		Version:           "6",
		ChainId:           chainId, // Mainnet
		VerifyingContract: "0x111111125421ca6dc452d289314280a0f8842a65",
	}

	// Convert order_data to Order struct with appropriate types
	salt := new(big.Int)
	salt.SetString(order_data["salt"], 10)

	makingAmount := new(big.Int)
	makingAmount.SetString(order_data["makingAmount"], 10)

	takingAmount := new(big.Int)
	takingAmount.SetString(order_data["takingAmount"], 10)

	saltBigInt, success := new(big.Int).SetString(order_data["salt"], 10)
	if !success {
		fmt.Println("error converting salt to big int")
		return
	}
	makingAmountBigInt, success := new(big.Int).SetString(order_data["makingAmount"], 10)
	if !success {
		fmt.Println("error converting makingAmount to big int")
		return
	}
	takingAmountBigInt, success := new(big.Int).SetString(order_data["takingAmount"], 10)
	if !success {
		fmt.Println("error converting takingAmount to big int")
		return
	}
	makerTraitsBigInt, success := new(big.Int).SetString(order_data["makerTraits"], 10)
	if !success {
		fmt.Println("error converting makerTraits to big int")
		return
	}

	orderMessage := apitypes.TypedDataMessage{
		"salt":         saltBigInt,
		"maker":        order_data["maker"],
		"receiver":     order_data["receiver"],
		"makerAsset":   order_data["makerAsset"],
		"takerAsset":   order_data["takerAsset"],
		"makingAmount": makingAmountBigInt,
		"takingAmount": takingAmountBigInt,
		"makerTraits":  makerTraitsBigInt,
	}

	typedData := apitypes.TypedData{
		Types: map[string][]apitypes.Type{
			"EIP712Domain": {
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"Order": {
				{Name: "salt", Type: "uint256"},
				{Name: "maker", Type: "address"},
				{Name: "receiver", Type: "address"},
				{Name: "makerAsset", Type: "address"},
				{Name: "takerAsset", Type: "address"},
				{Name: "makingAmount", Type: "uint256"},
				{Name: "takingAmount", Type: "uint256"},
				{Name: "makerTraits", Type: "uint256"},
			},
		},
		PrimaryType: "Order",
		Domain: apitypes.TypedDataDomain{
			Name:              domainData.Name,
			Version:           domainData.Version,
			ChainId:           domainData.ChainId,
			VerifyingContract: domainData.VerifyingContract,
		},
		Message: orderMessage,
	}

	//print out the typed data
	typedDataHash, _ := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	domainSeparator, _ := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())

	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	challengeHash := crypto.Keccak256Hash(rawData)
	challengeHashHex := challengeHash.Hex()
	fmt.Println("challengeHash:", challengeHashHex)

	privateKey, err := crypto.HexToECDSA("965e092fdfc08940d2bd05c7b5c7e1c51e283e92c7f52bbf1408973ae9a9acb7")
	if err != nil {
		fmt.Println("error converting private key to ECDSA:", err)
		return
	}

	// Sign the challenge hash
	signature, err := crypto.Sign(challengeHash.Bytes(), privateKey)
	if err != nil {
		fmt.Println("error signing challenge hash:", err)
		return
	}

	// add 27 to `v` value (last byte) because reasons I don't care to understand
	signature[64] += 27

	// convert signature to hex string
	signatureHex := fmt.Sprintf("0x%x", signature)
	fmt.Println("signature:", signatureHex)

	// Construct the body
	body := map[string]interface{}{
		// for `data` we need to convert salt to an int but the rest need to be strings
		"order":     order_data,
		"signature": signatureHex,
		// "orderHash": challengeHashHex,
		"extension": extension,
		"quoteId":   "string", // this is a placeholder, replace with actual quoteId if needed
	}

	// Convert the body to JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	// Print the JSON body to be sent
	fmt.Println("JSON Body to be sent:", string(jsonBody)) // fine here

	// Define the URL
	// convert chainId to string
	url := "https://api.1inch.dev/orderbook/v4.0/" + fmt.Sprintf("%d", chainIdInt)

	// Create a new request using http
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	// Read and print the request body
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("Error reading request body: %v\n", err)
		return
	}
	// Replace the body for future use
	req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

	// Now print the body
	// fmt.Println("Request Body:", string(reqBody)) // This should print the actual body content

	// Add the required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("accept", "application/json, text/plain, */*")

	// Send the request via a client
	// Send the request via a client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if buf.Len() == 0 {
		log.Fatal("Empty response")
	}

	fmt.Printf("\n\nresponse Status: %s", buf)

}
