package tokenprices

type PricesParameters struct {
	Currency CurrencyType
}

type CurrencyType string

const (
	CurrencyTypeWEI CurrencyType = ""
	CurrencyTypeUSD              = "USD"
	CurrencyTypeAED              = "AED"
	CurrencyTypeARS              = "ARS"
	CurrencyTypeAUD              = "AUD"
	CurrencyTypeBDT              = "BDT"
	CurrencyTypeBHD              = "BHD"
	CurrencyTypeBMD              = "BMD"
	CurrencyTypeBRL              = "BRL"
	CurrencyTypeCAD              = "CAD"
	CurrencyTypeCHF              = "CHF"
	CurrencyTypeCLP              = "CLP"
	CurrencyTypeCNY              = "CNY"
	CurrencyTypeCZK              = "CZK"
	CurrencyTypeDKK              = "DKK"
	CurrencyTypeEUR              = "EUR"
	CurrencyTypeGBP              = "GBP"
	CurrencyTypeHKD              = "HKD"
	CurrencyTypeHUF              = "HUF"
	CurrencyTypeIDR              = "IDR"
	CurrencyTypeILS              = "ILS"
	CurrencyTypeINR              = "INR"
	CurrencyTypeJPY              = "JPY"
	CurrencyTypeKRW              = "KRW"
	CurrencyTypeKWD              = "KWD"
	CurrencyTypeLKR              = "LKR"
	CurrencyTypeMMK              = "MMK"
	CurrencyTypeMXN              = "MXN"
	CurrencyTypeMYR              = "MYR"
	CurrencyTypeNGN              = "NGN"
	CurrencyTypeNOK              = "NOK"
	CurrencyTypeNZD              = "NZD"
	CurrencyTypePHP              = "PHP"
	CurrencyTypePKR              = "PKR"
	CurrencyTypePLN              = "PLN"
	CurrencyTypeRUB              = "RUB"
	CurrencyTypeSAR              = "SAR"
	CurrencyTypeSEK              = "SEK"
	CurrencyTypeSGD              = "SGD"
	CurrencyTypeTHB              = "THB"
	CurrencyTypeTRY              = "TRY"
	CurrencyTypeTWD              = "TWD"
	CurrencyTypeUAH              = "UAH"
	CurrencyTypeVEF              = "VEF"
	CurrencyTypeVND              = "VND"
	CurrencyTypeZAR              = "ZAR"
)

var CurrencyTypeValues = []CurrencyType{
	CurrencyTypeWEI,
	CurrencyTypeUSD,
	CurrencyTypeAED,
	CurrencyTypeARS,
	CurrencyTypeAUD,
	CurrencyTypeBDT,
	CurrencyTypeBHD,
	CurrencyTypeBMD,
	CurrencyTypeBRL,
	CurrencyTypeCAD,
	CurrencyTypeCHF,
	CurrencyTypeCLP,
	CurrencyTypeCNY,
	CurrencyTypeCZK,
	CurrencyTypeDKK,
	CurrencyTypeEUR,
	CurrencyTypeGBP,
	CurrencyTypeHKD,
	CurrencyTypeHUF,
	CurrencyTypeIDR,
	CurrencyTypeILS,
	CurrencyTypeINR,
	CurrencyTypeJPY,
	CurrencyTypeKRW,
	CurrencyTypeKWD,
	CurrencyTypeLKR,
	CurrencyTypeMMK,
	CurrencyTypeMXN,
	CurrencyTypeMYR,
	CurrencyTypeNGN,
	CurrencyTypeNOK,
	CurrencyTypeNZD,
	CurrencyTypePHP,
	CurrencyTypePKR,
	CurrencyTypePLN,
	CurrencyTypeRUB,
	CurrencyTypeSAR,
	CurrencyTypeSEK,
	CurrencyTypeSGD,
	CurrencyTypeTHB,
	CurrencyTypeTRY,
	CurrencyTypeTWD,
	CurrencyTypeUAH,
	CurrencyTypeVEF,
	CurrencyTypeVND,
	CurrencyTypeZAR,
}
