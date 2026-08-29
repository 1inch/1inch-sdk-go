package traces

// Clean names matching the traces methods; the generated Dto-suffixed names and
// the singular param name remain valid.
type (
	GetSyncedIntervalResponse     = ReadSyncedIntervalResponseDto
	GetBlockTraceByNumberResponse = CoreBuiltinBlockTracesDto
	GetBlockTraceByNumberParams   = GetBlockTraceByNumberParam
)
