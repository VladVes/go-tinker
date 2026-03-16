package schemas 

type (
	CreateLogEntryRequest struct {
		Message   string `json:"message"`
		Level     string `json:"level"`
		Timestamp int64  `json:"timestamp"`
	}

	CreateLogEntryResponse struct {
		ID string `json:"id"`
	}

	LogEntry struct {
		ID        string
		Message   string
		Level     string
		Timestamp int64
	}
)

