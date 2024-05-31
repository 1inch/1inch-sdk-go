package portfolio

//TODO These response bodies are largely a guess due to the lack of documentation

type GetPortfolioValueResponse struct {
	Result []struct {
		ProtocolName   string `json:"protocol_name"`
		ProtocolLocked bool   `json:"protocol_locked"`
		Result         []struct {
			ChainId  *int    `json:"chain_id"`
			ValueUsd float64 `json:"value_usd"`
		} `json:"result"`
	} `json:"result"`
	System struct {
		ClickTime         float64 `json:"click_time"`
		NodeTime          float64 `json:"node_time"`
		MicroservicesTime float64 `json:"microservices_time"`
		RedisTime         float64 `json:"redis_time"`
		TotalTime         float64 `json:"total_time"`
	} `json:"system"`
}

type GetPortfolioProfitAndLossResponse struct {
	Result []struct {
		ChainId      *int    `json:"chain_id"`
		AbsProfitUsd float64 `json:"abs_profit_usd"`
		Roi          float64 `json:"roi"`
	} `json:"result"`
	System struct {
		ClickTime         float64 `json:"click_time"`
		NodeTime          float64 `json:"node_time"`
		MicroservicesTime float64 `json:"microservices_time"`
		RedisTime         float64 `json:"redis_time"`
		TotalTime         float64 `json:"total_time"`
	} `json:"system"`
}

type GetProtocolsDetailsResponse struct {
	Result []struct {
		ChainId          int      `json:"chain_id"`
		ContractAddress  string   `json:"contract_address"`
		Addresses        []string `json:"addresses"`
		TokenId          int      `json:"token_id"`
		BlockNumber      int      `json:"block_number"`
		Protocol         string   `json:"protocol"`
		Name             string   `json:"name"`
		ContractType     string   `json:"contract_type"`
		SubContractType  string   `json:"sub_contract_type"`
		IsWhitelisted    int      `json:"is_whitelisted"`
		Status           int      `json:"status"`
		UnderlyingTokens []struct {
			Address    string  `json:"address"`
			Decimals   int     `json:"decimals"`
			Amount     float64 `json:"amount"`
			PriceToUsd float64 `json:"price_to_usd"`
			ValueUsd   float64 `json:"value_usd"`
		} `json:"underlying_tokens"`
		ValueUsd     float64 `json:"value_usd"`
		ProtocolName string  `json:"protocol_name"`
		Info         struct {
			Amount                    float64     `json:"amount"`
			UnderlyingContractAddress float64     `json:"underlying_contract_address"`
			PositionStartTimestamp    interface{} `json:"position_start_timestamp"`
			TimeIntervalSec           int         `json:"time_interval_sec,omitempty"`
			PriceToUsd                float64     `json:"price_to_usd"`
			ValueUsd                  float64     `json:"value_usd"`
			GainedReward              float64     `json:"gained_reward"`
			GainedRewardUsd           float64     `json:"gained_reward_usd"`
			Roi                       *float64    `json:"roi"`
			Apr                       *float64    `json:"apr"`
			Apy                       *float64    `json:"apy,omitempty"`
			AveragePriceUsd           *float64    `json:"average_price_usd,omitempty"`
			HoldingTimeDays           float64     `json:"holding_time_days"`
			StethAmount               float64     `json:"steth_amount,omitempty"`
			IsWrapped                 int         `json:"is_wrapped,omitempty"`
			ChiCoeff                  float64     `json:"chi_coeff,omitempty"`
			ChaiAmount                float64     `json:"chai_amount,omitempty"`
			UnderlyingAmountInvested  float64     `json:"underlying_amount_invested,omitempty"`
			UnderlyingAmount          float64     `json:"underlying_amount,omitempty"`
		} `json:"info"`
	} `json:"result"`
	System struct {
		ClickTime         float64 `json:"click_time"`
		NodeTime          float64 `json:"node_time"`
		MicroservicesTime float64 `json:"microservices_time"`
		RedisTime         float64 `json:"redis_time"`
		TotalTime         float64 `json:"total_time"`
	} `json:"system"`
}

type GetTokensCurrentValueResponse struct {
	Result []struct {
		ProtocolName   string `json:"protocol_name"`
		ProtocolLocked bool   `json:"protocol_locked"`
		Result         []struct {
			ChainId  *int    `json:"chain_id"`
			ValueUsd float64 `json:"value_usd"`
		} `json:"result"`
	} `json:"result"`
	System struct {
		ClickTime         float64 `json:"click_time"`
		NodeTime          float64 `json:"node_time"`
		MicroservicesTime float64 `json:"microservices_time"`
		RedisTime         float64 `json:"redis_time"`
		TotalTime         float64 `json:"total_time"`
	} `json:"system"`
}

type GetTokensProfitLossResponse struct {
	Result []struct {
		ChainId      *int    `json:"chain_id"`
		AbsProfitUsd float64 `json:"abs_profit_usd"`
		Roi          float64 `json:"roi"`
	} `json:"result"`
	System struct {
		ClickTime         float64 `json:"click_time"`
		NodeTime          float64 `json:"node_time"`
		MicroservicesTime float64 `json:"microservices_time"`
		RedisTime         float64 `json:"redis_time"`
		TotalTime         float64 `json:"total_time"`
	} `json:"system"`
}

type GetTokensDetailsResponse struct {
	Result []struct {
		ChainId         int      `json:"chain_id"`
		ContractAddress string   `json:"contract_address"`
		PriceToUsd      float64  `json:"price_to_usd"`
		Amount          float64  `json:"amount"`
		ValueUsd        float64  `json:"value_usd"`
		AbsProfitUsd    *float64 `json:"abs_profit_usd"`
		Roi             *float64 `json:"roi"`
		Status          int      `json:"status"`
	} `json:"result"`
	System struct {
		ClickTime         float64 `json:"click_time"`
		NodeTime          float64 `json:"node_time"`
		MicroservicesTime float64 `json:"microservices_time"`
		RedisTime         float64 `json:"redis_time"`
		TotalTime         float64 `json:"total_time"`
	} `json:"system"`
}

type IsServiceAvailableResponse struct {
	Result bool `json:"result"`
}

type GetSupportedChainsResponse struct {
	Result []int `json:"result"`
}

type GetCurrentValueResponse struct {
	Result []struct {
		Address  string  `json:"address"`
		ValueUsd float64 `json:"value_usd"`
	} `json:"result"`
	System struct {
		ClickTime         float64 `json:"click_time"`
		NodeTime          float64 `json:"node_time"`
		MicroservicesTime float64 `json:"microservices_time"`
		RedisTime         float64 `json:"redis_time"`
		TotalTime         float64 `json:"total_time"`
	} `json:"system"`
}

type GetCurrentProfitLossResponse struct {
	Result []struct {
		ChainId      *int    `json:"chain_id"`
		AbsProfitUsd float64 `json:"abs_profit_usd"`
		Roi          float64 `json:"roi"`
	} `json:"result"`
	System struct {
		ClickTime         float64 `json:"click_time"`
		NodeTime          float64 `json:"node_time"`
		MicroservicesTime float64 `json:"microservices_time"`
		RedisTime         float64 `json:"redis_time"`
		TotalTime         float64 `json:"total_time"`
	} `json:"system"`
}

type GetValueChartResponse struct {
	Result []struct {
		Timestamp int     `json:"timestamp"`
		ValueUsd  float64 `json:"value_usd"`
	} `json:"result"`
	System struct {
		ClickTime         float64 `json:"click_time"`
		NodeTime          float64 `json:"node_time"`
		MicroservicesTime float64 `json:"microservices_time"`
		RedisTime         float64 `json:"redis_time"`
		TotalTime         float64 `json:"total_time"`
	} `json:"system"`
}
