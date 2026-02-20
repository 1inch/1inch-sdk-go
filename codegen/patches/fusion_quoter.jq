# Fix PresetClass.exclusiveResolver: "object" -> "string"
# The API returns a string address, not a JSON object.
.components.schemas.PresetClass.properties.exclusiveResolver.type = "string"
