package history

var (
	_ Item    = HistoryEvent{}
	_ Details = HistoryEventDetails{}
)
