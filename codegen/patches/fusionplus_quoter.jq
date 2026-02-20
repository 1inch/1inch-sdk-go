# Fix GetQuoteOutput.quoteId: "object" -> "string"
# The API returns a string, not a JSON object.
.components.schemas.GetQuoteOutput.properties.quoteId.type = "string"

# Fix Preset.exclusiveResolver: "object" -> "string"
| .components.schemas.Preset.properties.exclusiveResolver.type = "string"

# Fix amount parameter: "number" -> "string" (BigInt values overflow float)
| .paths["/v1.0/quote/receive"].post.parameters |= map(
    if .name == "amount" then .schema.type = "string" else . end
  )

# Fix fee parameter: "number" -> "string" (BigInt values overflow float)
| .paths["/v1.0/quote/receive"].post.parameters |= map(
    if .name == "fee" then .schema.type = "string" else . end
  )

# Fix isPermit2 parameter: "string" -> "boolean"
| .paths["/v1.0/quote/receive"].post.parameters |= map(
    if .name == "isPermit2" then .schema.type = "boolean" else . end
  )
