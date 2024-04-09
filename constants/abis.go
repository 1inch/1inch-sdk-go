package constants

import _ "embed"

//go:embed abi/erc20.abi.json
var Erc20ABI string

//go:embed abi/seriesNonceManager.abi.json
var SeriesNonceManagerABI string

//go:embed abi/aggregationRouterV5.abi.json
var AggregationRouterV5ABI string
