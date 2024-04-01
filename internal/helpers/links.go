package helpers

import (
	"fmt"

	"github.com/1inch/1inch-sdk-go/constants"
)

//func PrintBlockExplorerTxLink(chainId int, txHash string) {
//	output := GetBlockExplorerTxLinkInfo(chainId, txHash)
//	fmt.Printf(output)
//}

func GetBlockExplorerTxLinkInfo(chainId int, txHash string) string {
	const etherscanBaseURL = "https://etherscan.io/tx/"
	const etherscan = "Etherscan"
	const polygonScanBaseURL = "https://polygonscan.com/tx/"
	const polygonScan = "PolygonScan"

	var baseUrl, serviceName string
	switch chainId {
	case constants.Ethereum:
		baseUrl = etherscanBaseURL
		serviceName = etherscan
	case constants.Polygon:
		baseUrl = polygonScanBaseURL
		serviceName = polygonScan
	default:
		return fmt.Sprintf("Tx Id: %s\n", txHash)

	}
	return fmt.Sprintf("View it on %s here: %s\n", serviceName, baseUrl+txHash)
}
