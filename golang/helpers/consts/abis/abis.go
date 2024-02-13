package abis

import _ "embed"

//go:embed erc20.abi.json
var Erc20 string

//go:embed seriesNonceManager.abi.json
var SeriesNonceManager string

//go:embed aggregationRouterV5.abi.json
var AggregationRouterV5 string
