# Fix GetOrderFillsByHashOutput.dstTokenPriceUsd: "object" -> "string"
.components.schemas.GetOrderFillsByHashOutput.properties.dstTokenPriceUsd.type = "string"

# Fix GetOrderFillsByHashOutput.srcTokenPriceUsd: "object" -> "string"
| .components.schemas.GetOrderFillsByHashOutput.properties.srcTokenPriceUsd.type = "string"

# Fix GetOrderFillsByHashOutput.points: single $ref -> array of $ref
| .components.schemas.GetOrderFillsByHashOutput.properties.points = {
    "type": "array",
    "items": { "$ref": "#/components/schemas/AuctionPointOutput" },
    "x-go-type-skip-optional-pointer": true
  }
