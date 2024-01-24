package tenderly

import "time"

type SwapConfig struct {
	ChainId         int
	PublicAddress   string
	FromToken       string
	Value           string
	TransactionData string
	ApproveFirst    bool
}

type RunConfiguration struct {
	Environment           string `yaml:"environment"`           // Determines the Pathfinder and Swapbuilder services to use
	ChainId               string `yaml:"chainId"`               // The ID of the chainId to use (ex. Ethereum would be chain-id=1)
	DexId                 string `yaml:"dexId"`                 // The ID of the DEX to test (ex. dex-id=INTEGRAL)
	FromToken             string `yaml:"fromToken"`             // The token address to swap from (source token) - will be lowercased automatically
	ToToken               string `yaml:"toToken"`               // The token address to swap to (destination token) - will be lowercased automatically
	Amount                string `yaml:"amount"`                // Amount of from-token being swapped (ex. 1 ETH would be 1000000000000000000)
	AmountEth             string `yaml:"-"`                     // Populated dynamically if FromToken is ETH
	FromWallet            string `yaml:"fromWallet"`            // The wallet that will be the source of the transaction
	OpenBrowserOnComplete bool   `yaml:"openBrowserOnComplete"` // When set to true, automatically opens the browser with the results of the Tenderly simulation
	PathfinderVersion     string `yaml:"pathfinderVersion"`     // The endpoint version of Pathfinder requests (only v1.2 is officially supported)
	SwapbuilderVersion    string `yaml:"swapbuilderVersion"`    // The endpoint version of Swapbuilder requests (only v5.0 is officially supported)
}

type InputData struct {
	Contract string
	Calldata string
	GasLimit int
	Root     string
}

type SimulateRequest struct {
	From               string                 `json:"from"`
	To                 string                 `json:"to"`
	Input              string                 `json:"input"`
	Gas                int                    `json:"gas"`
	GasPrice           string                 `json:"gas_price"`
	Value              string                 `json:"value"`
	Save               bool                   `json:"save"`
	GenerateAccessList bool                   `json:"generate_access_list"`
	SaveIfFails        bool                   `json:"save_if_fails"`
	Root               string                 `json:"root,omitempty"`
	SimulationType     string                 `json:"simulation_type"`
	StateObjects       map[string]StateObject `json:"state_objects"`
}

type StateObject struct {
	Balance string            `json:"balance,omitempty"`
	Storage map[string]string `json:"storage,omitempty"`
}

type ForkRequest struct {
	NetworkID string `json:"network_id"`
	ForkName  string `json:"alias"`
}

type ForkResponse struct {
	SimulationFork struct {
		ForkID string `json:"id"`
	} `json:"simulation_fork"`
}

type ResponseError struct {
	ErrorStruct struct {
		Message string `json:"message"`
		Slug    string `json:"slug"`
	} `json:"error"`
}

func (te ResponseError) Error() string {
	return te.ErrorStruct.Message
}

type SimulationResponse struct {
	Transaction struct {
		Hash              string      `json:"hash"`
		BlockHash         string      `json:"block_hash"`
		BlockNumber       int         `json:"block_number"`
		From              string      `json:"from"`
		Gas               int         `json:"gas"`
		GasPrice          int64       `json:"gas_price"`
		GasFeeCap         int         `json:"gas_fee_cap"`
		GasTipCap         int         `json:"gas_tip_cap"`
		CumulativeGasUsed int         `json:"cumulative_gas_used"`
		GasUsed           int         `json:"gas_used"`
		EffectiveGasPrice int         `json:"effective_gas_price"`
		Input             string      `json:"input"`
		Nonce             int         `json:"nonce"`
		To                string      `json:"to"`
		Index             int         `json:"index"`
		Value             string      `json:"value"`
		AccessList        interface{} `json:"access_list"`
		Status            bool        `json:"status"`
		Addresses         interface{} `json:"addresses"`
		ContractIds       interface{} `json:"contract_ids"`
		NetworkID         string      `json:"network_id"`
		Timestamp         time.Time   `json:"timestamp"`
		FunctionSelector  string      `json:"function_selector"`
		L1BlockNumber     int         `json:"l1_block_number"`
		L1Timestamp       int         `json:"l1_timestamp"`
		DepositTx         bool        `json:"deposit_tx"`
		Mint              interface{} `json:"mint"`
		Sig               struct {
			V string `json:"v"`
			R string `json:"r"`
			S string `json:"s"`
		} `json:"sig"`
		TransactionInfo struct {
			ContractID      string      `json:"contract_id"`
			BlockNumber     int         `json:"block_number"`
			TransactionID   string      `json:"transaction_id"`
			ContractAddress string      `json:"contract_address"`
			Method          string      `json:"method"`
			Parameters      interface{} `json:"parameters"`
			IntrinsicGas    int         `json:"intrinsic_gas"`
			RefundGas       int         `json:"refund_gas"`
			CallTrace       struct {
				Hash               string `json:"hash"`
				ContractName       string `json:"contract_name"`
				FunctionName       string `json:"function_name"`
				FunctionPc         int    `json:"function_pc"`
				FunctionOp         string `json:"function_op"`
				FunctionFileIndex  int    `json:"function_file_index"`
				FunctionCodeStart  int    `json:"function_code_start"`
				FunctionLineNumber int    `json:"function_line_number"`
				FunctionCodeLength int    `json:"function_code_length"`
				AbsolutePosition   int    `json:"absolute_position"`
				CallerPc           int    `json:"caller_pc"`
				CallerOp           string `json:"caller_op"`
				CallType           string `json:"call_type"`
				From               string `json:"from"`
				FromBalance        string `json:"from_balance"`
				To                 string `json:"to"`
				ToBalance          string `json:"to_balance"`
				Value              string `json:"value"`
				Caller             struct {
					Address string `json:"address"`
					Balance string `json:"balance"`
				} `json:"caller"`
				BlockTimestamp time.Time `json:"block_timestamp"`
				Gas            int       `json:"gas"`
				GasUsed        int       `json:"gas_used"`
				IntrinsicGas   int       `json:"intrinsic_gas"`
				Input          string    `json:"input"`
				DecodedInput   []struct {
					Soltype struct {
						Name            string      `json:"name"`
						Type            string      `json:"type"`
						StorageLocation string      `json:"storage_location"`
						Components      interface{} `json:"components"`
						Offset          int         `json:"offset"`
						Index           string      `json:"index"`
						Indexed         bool        `json:"indexed"`
						SimpleType      struct {
							Type string `json:"type"`
						} `json:"simple_type"`
					} `json:"soltype,omitempty"`
					Value    string `json:"value"`
					Soltype0 struct {
						Name            string `json:"name"`
						Type            string `json:"type"`
						StorageLocation string `json:"storage_location"`
						Components      []struct {
							Name            string      `json:"name"`
							Type            string      `json:"type"`
							StorageLocation string      `json:"storage_location"`
							Components      interface{} `json:"components"`
							Offset          int         `json:"offset"`
							Index           string      `json:"index"`
							Indexed         bool        `json:"indexed"`
							SimpleType      struct {
								Type string `json:"type"`
							} `json:"simple_type"`
						} `json:"components"`
						Offset  int    `json:"offset"`
						Index   string `json:"index"`
						Indexed bool   `json:"indexed"`
					} `json:"soltype,omitempty"`
				} `json:"decoded_input"`
				BalanceDiff []struct {
					Address  string `json:"address"`
					Original string `json:"original"`
					Dirty    string `json:"dirty"`
					IsMiner  bool   `json:"is_miner"`
				} `json:"balance_diff"`
				NonceDiff []struct {
					Address  string `json:"address"`
					Original string `json:"original"`
					Dirty    string `json:"dirty"`
				} `json:"nonce_diff"`
				Output          string      `json:"output"`
				DecodedOutput   interface{} `json:"decoded_output"`
				Error           string      `json:"error"`
				ErrorOp         string      `json:"error_op"`
				ErrorFileIndex  int         `json:"error_file_index"`
				ErrorLineNumber int         `json:"error_line_number"`
				ErrorCodeStart  int         `json:"error_code_start"`
				ErrorCodeLength int         `json:"error_code_length"`
				NetworkID       string      `json:"network_id"`
				Calls           []struct {
					Hash               string      `json:"hash"`
					ContractName       string      `json:"contract_name"`
					FunctionName       string      `json:"function_name"`
					FunctionPc         int         `json:"function_pc"`
					FunctionOp         string      `json:"function_op"`
					FunctionFileIndex  int         `json:"function_file_index,omitempty"`
					FunctionCodeStart  int         `json:"function_code_start,omitempty"`
					FunctionLineNumber int         `json:"function_line_number,omitempty"`
					FunctionCodeLength int         `json:"function_code_length,omitempty"`
					AbsolutePosition   int         `json:"absolute_position"`
					CallerPc           int         `json:"caller_pc"`
					CallerOp           string      `json:"caller_op"`
					CallerFileIndex    int         `json:"caller_file_index,omitempty"`
					CallerLineNumber   int         `json:"caller_line_number,omitempty"`
					CallerCodeStart    int         `json:"caller_code_start,omitempty"`
					CallerCodeLength   int         `json:"caller_code_length,omitempty"`
					CallType           string      `json:"call_type,omitempty"`
					From               string      `json:"from"`
					FromBalance        interface{} `json:"from_balance"`
					To                 string      `json:"to"`
					ToBalance          interface{} `json:"to_balance"`
					Value              interface{} `json:"value"`
					Caller             struct {
						Address string `json:"address"`
						Balance string `json:"balance"`
					} `json:"caller,omitempty"`
					BlockTimestamp time.Time `json:"block_timestamp"`
					Gas            int       `json:"gas"`
					GasUsed        int       `json:"gas_used"`
					Input          string    `json:"input"`
					DecodedInput   []struct {
						Soltype struct {
							Name            string      `json:"name"`
							Type            string      `json:"type"`
							StorageLocation string      `json:"storage_location"`
							Components      interface{} `json:"components"`
							Offset          int         `json:"offset"`
							Index           string      `json:"index"`
							Indexed         bool        `json:"indexed"`
							SimpleType      struct {
								Type string `json:"type"`
							} `json:"simple_type"`
						} `json:"soltype"`
						Value string `json:"value"`
					} `json:"decoded_input,omitempty"`
					Output        string `json:"output"`
					DecodedOutput []struct {
						Soltype struct {
							Name            string      `json:"name"`
							Type            string      `json:"type"`
							StorageLocation string      `json:"storage_location"`
							Components      interface{} `json:"components"`
							Offset          int         `json:"offset"`
							Index           string      `json:"index"`
							Indexed         bool        `json:"indexed"`
							SimpleType      struct {
								Type string `json:"type"`
							} `json:"simple_type"`
						} `json:"soltype"`
						Value bool `json:"value"`
					} `json:"decoded_output"`
					NetworkID       string      `json:"network_id"`
					Calls           interface{} `json:"calls"`
					Error           string      `json:"error,omitempty"`
					ErrorOp         string      `json:"error_op,omitempty"`
					ErrorFileIndex  int         `json:"error_file_index,omitempty"`
					ErrorLineNumber int         `json:"error_line_number,omitempty"`
					ErrorCodeStart  int         `json:"error_code_start,omitempty"`
					ErrorCodeLength int         `json:"error_code_length,omitempty"`
				} `json:"calls"`
			} `json:"call_trace"`
			StackTrace []struct {
				FileIndex int    `json:"file_index"`
				Contract  string `json:"contract"`
				Name      string `json:"name"`
				Line      int    `json:"line"`
				Error     string `json:"error"`
				Code      string `json:"code"`
				Op        string `json:"op"`
				Length    int    `json:"length"`
			} `json:"stack_trace"`
			Logs        interface{} `json:"logs"`
			BalanceDiff []struct {
				Address  string `json:"address"`
				Original string `json:"original"`
				Dirty    string `json:"dirty"`
				IsMiner  bool   `json:"is_miner"`
			} `json:"balance_diff"`
			NonceDiff []struct {
				Address  string `json:"address"`
				Original string `json:"original"`
				Dirty    string `json:"dirty"`
			} `json:"nonce_diff"`
			StateDiff    interface{} `json:"state_diff"`
			RawStateDiff interface{} `json:"raw_state_diff"`
			ConsoleLogs  interface{} `json:"console_logs"`
			CreatedAt    time.Time   `json:"created_at"`
		} `json:"transaction_info"`
		ErrorMessage string `json:"error_message"`
		ErrorInfo    struct {
			ErrorMessage string `json:"error_message"`
			Address      string `json:"address"`
		} `json:"error_info"`
		Method       string      `json:"method"`
		DecodedInput interface{} `json:"decoded_input"`
		CallTrace    []struct {
			CallType    string `json:"call_type"`
			From        string `json:"from"`
			To          string `json:"to"`
			Gas         int    `json:"gas"`
			GasUsed     int    `json:"gas_used"`
			Value       string `json:"value"`
			Error       string `json:"error"`
			Type        string `json:"type"`
			Input       string `json:"input"`
			FromBalance string `json:"fromBalance"`
			ToBalance   string `json:"toBalance"`
		} `json:"call_trace"`
	} `json:"transaction"`
	Simulation struct {
		ID               string      `json:"id"`
		ProjectID        string      `json:"project_id"`
		OwnerID          string      `json:"owner_id"`
		NetworkID        string      `json:"network_id"`
		BlockNumber      int         `json:"block_number"`
		TransactionIndex int         `json:"transaction_index"`
		From             string      `json:"from"`
		To               string      `json:"to"`
		Input            string      `json:"input"`
		Gas              int         `json:"gas"`
		GasPrice         string      `json:"gas_price"`
		Value            string      `json:"value"`
		Status           bool        `json:"status"`
		AccessList       interface{} `json:"access_list"`
		QueueOrigin      string      `json:"queue_origin"`
		BlockHeader      struct {
			Number           string      `json:"number"`
			Hash             string      `json:"hash"`
			StateRoot        string      `json:"stateRoot"`
			ParentHash       string      `json:"parentHash"`
			Sha3Uncles       string      `json:"sha3Uncles"`
			TransactionsRoot string      `json:"transactionsRoot"`
			ReceiptsRoot     string      `json:"receiptsRoot"`
			LogsBloom        string      `json:"logsBloom"`
			Timestamp        string      `json:"timestamp"`
			Difficulty       string      `json:"difficulty"`
			GasLimit         string      `json:"gasLimit"`
			GasUsed          string      `json:"gasUsed"`
			Miner            string      `json:"miner"`
			ExtraData        string      `json:"extraData"`
			MixHash          string      `json:"mixHash"`
			Nonce            string      `json:"nonce"`
			BaseFeePerGas    string      `json:"baseFeePerGas"`
			Size             string      `json:"size"`
			TotalDifficulty  string      `json:"totalDifficulty"`
			Uncles           interface{} `json:"uncles"`
			Transactions     interface{} `json:"transactions"`
		} `json:"block_header"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"simulation"`
	Contracts           []interface{} `json:"contracts"`
	GeneratedAccessList []struct {
		Address string `json:"address"`
	} `json:"generated_access_list"`
}
