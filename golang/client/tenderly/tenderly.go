package tenderly

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/1inch/1inch-sdk/golang/helpers/consts/contracts"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/tokens"
)

func SimulateSwap(tenderlyApiKey string, config SwapConfig) (*SimulationResponse, error) {
	name := "Dev Portal - Swap"
	if config.ApproveFirst {
		name += " with approval"
	}
	forkId, err := CreateTenderlyFork(tenderlyApiKey, config.ChainId, name)
	if err != nil {
		return nil, fmt.Errorf("failed to create tenderly fork: %v", err)
	}

	var root string
	if config.ApproveFirst {

		fmt.Printf("Tenderly: Non-ETH token is being swapped! Executing request to approve %s for swapping\n", config.FromToken)

		// Static calldata that approves a large ERC20 spend limit for the v5 router
		const ApproveErc20Calldata = `0x095ea7b30000000000000000000000001111111254eeb25477b68fb85ed929f73a960582000000000000000000000000000000000000000c9f2c9cd04674edea3fffffff`
		const ApproveErc20GasLimitStatic = 2000000

		fmt.Println("Tenderly: Simulating token approval on Tenderly")

		tokenApprovalSimulationRequest := &SimulateRequest{
			From:               config.PublicAddress,
			To:                 config.FromToken,
			Input:              ApproveErc20Calldata,
			Gas:                ApproveErc20GasLimitStatic,
			GasPrice:           "1806564247",
			Value:              "0",
			Save:               true,
			GenerateAccessList: true,
			SaveIfFails:        true,
			SimulationType:     "quick",
		}
		tokenApprovalResponse, errApprove := ExecuteTenderlySimulationRequest(tenderlyApiKey, forkId, tokenApprovalSimulationRequest)
		if errApprove != nil {
			return nil, fmt.Errorf("request to approve tokens on Tenderly failed: %v\n", errApprove)
		}

		root = tokenApprovalResponse.Simulation.ID
		preApprovalSimulationUrl := fmt.Sprintf("https://dashboard.tenderly.co/Natalia/backend-/fork/%s/simulation/%s", forkId, root)
		fmt.Printf("Tenderly: Pre-approval simulation complete. Link to results: %s\n", preApprovalSimulationUrl)
	}

	// Always set the balance of the sending wallet to 1 ETH
	stateOverrides := map[string]StateObject{
		config.PublicAddress: {
			Balance: "1000000000000000000", // 1 Ether in wei
		},
	}

	// Temporary hard coding 1INCH swaps to use a static balance of 100 to save the implementation logic
	if config.FromToken == tokens.Ethereum1inch {
		stateOverrides["0x111111111117dc0aa78b770fa6a738034120c302"] = StateObject{
			Storage: map[string]string{
				"0x39c953eb19bbd67c4c3fd8344cb9af0980b96fc08605b375e22a196c732d92a8": "0x0000000000000000000000000000000000000000000000056bc75e2d63100000",
			},
		}
	}

	swapSimulationResponse, err := ExecuteTenderlySimulationRequest(tenderlyApiKey, forkId, &SimulateRequest{
		From:               config.PublicAddress,
		To:                 contracts.AggregationRouterV5,
		Input:              config.TransactionData,
		Gas:                30000000, // picked randomly
		Root:               root,     // TODO verify this works
		GasPrice:           "14806564247",
		Value:              config.Value,
		Save:               true,
		GenerateAccessList: true,
		SaveIfFails:        true,
		SimulationType:     "quick",
		StateObjects:       stateOverrides,
	})
	if err != nil {
		return nil, fmt.Errorf("request to tenderly failed: %v\n", err)
	}

	simulationUrl := fmt.Sprintf("https://dashboard.tenderly.co/Natalia/backend-/fork/%s/simulation/%s", forkId, swapSimulationResponse.Simulation.ID)
	fmt.Printf("Tenderly: Simulation complete. Link to results: %s\n", simulationUrl)

	transactionError := swapSimulationResponse.Transaction.TransactionInfo.CallTrace.Error
	if transactionError != "" {
		return nil, fmt.Errorf("simulation failed: %s\n", transactionError)
	}

	return swapSimulationResponse, nil
}

func ExecuteTenderlySimulationRequest(tenderlyApiKey string, forkId string, request *SimulateRequest) (*SimulationResponse, error) {

	base, err := url.Parse("https://api.tenderly.co")
	if err != nil {
		return nil, err
	}

	base.Path += fmt.Sprintf("/api/v1/account/Natalia/project/backend-/fork/%s/simulate", forkId)

	requestMarshaled, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, base.String(), io.NopCloser(bytes.NewBuffer(requestMarshaled)))
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Access-Key", tenderlyApiKey)

	httpClient := http.DefaultClient
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		var apiErr ResponseError
		if err2 := json.Unmarshal(data, &apiErr); err2 != nil {
			return nil, errors.New(strings.TrimSpace(string(data)))
		}
		return nil, apiErr
	}

	var tenderlyResponse SimulationResponse
	err = json.Unmarshal(data, &tenderlyResponse)
	if err != nil {
		return nil, err
	}

	return &tenderlyResponse, nil
}

func CreateTenderlyFork(tenderlyApiKey string, chainId int, alias string) (string, error) {

	base, err := url.Parse("https://api.tenderly.co")
	if err != nil {
		return "", err
	}

	base.Path += "/api/v1/account/Natalia/project/backend-/fork"

	requestBody := ForkRequest{
		NetworkID: fmt.Sprintf("%d", chainId),
		ForkName:  alias,
	}

	requestMarshaled, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, base.String(), io.NopCloser(bytes.NewBuffer(requestMarshaled)))
	if err != nil {
		return "", err
	}

	req.Header.Add("X-Access-Key", tenderlyApiKey)

	httpClient := http.DefaultClient
	res, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	err = res.Body.Close()
	if err != nil {
		return "", err
	}

	if res.StatusCode >= 400 {
		var apiErr ResponseError
		if err2 := json.Unmarshal(data, &apiErr); err2 != nil {
			return "", errors.New(strings.TrimSpace(string(data)))
		}
		return "", apiErr
	}

	var tenderlyForkResponse ForkResponse
	err = json.Unmarshal(data, &tenderlyForkResponse)
	if err != nil {
		return "", err
	}

	return tenderlyForkResponse.SimulationFork.ForkID, nil
}

func GetTenderlyForks(tenderlyApiKey string) (*GetForksResponse, error) {

	base, err := url.Parse("https://api.tenderly.co")
	if err != nil {
		return nil, err
	}

	base.Path += fmt.Sprintf("/api/v1/account/Natalia/project/backend-/forks")

	// Set query parameters
	query := base.Query()
	query.Set("page", "1")
	query.Set("perPage", "100")
	base.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, base.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Access-Key", tenderlyApiKey)

	httpClient := http.DefaultClient
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	var response GetForksResponse
	if res.StatusCode == 201 {
		err = json.Unmarshal(data, &response)
		if err != nil {
			return nil, err
		}
	} else if res.StatusCode == 404 {
		return nil, nil
	} else if res.StatusCode >= 400 {
		var apiErr ResponseError
		if err2 := json.Unmarshal(data, &apiErr); err2 != nil {
			return nil, errors.New(strings.TrimSpace(string(data)))
		}
		return nil, apiErr
	}

	return &response, nil
}

func DeleteTenderlyFork(tenderlyApiKey string, forkId string) error {

	base, err := url.Parse("https://api.tenderly.co")
	if err != nil {
		return err
	}

	base.Path += fmt.Sprintf("/api/v1/account/Natalia/project/backend-/fork/%s", forkId)

	req, err := http.NewRequest(http.MethodDelete, base.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Add("X-Access-Key", tenderlyApiKey)

	httpClient := http.DefaultClient
	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = res.Body.Close()
	if err != nil {
		return err
	}

	if res.StatusCode == 404 {
		return nil
	} else if res.StatusCode >= 400 {
		var apiErr ResponseError
		if err2 := json.Unmarshal(data, &apiErr); err2 != nil {
			return errors.New(strings.TrimSpace(string(data)))
		}
		return apiErr
	}

	return nil
}
