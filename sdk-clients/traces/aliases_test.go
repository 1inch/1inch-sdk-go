package traces

var (
	_ GetSyncedIntervalResponse     = ReadSyncedIntervalResponseDto{}
	_ GetBlockTraceByNumberResponse = CoreBuiltinBlockTracesDto{}
)
