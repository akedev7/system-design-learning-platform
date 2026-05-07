package response

type Status string

const (
	StatusSuccess Status = "success"
	StatusError  Status = "error"
)

type Envelope struct {
	Status  Status      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func Success(data interface{}) Envelope {
	return Envelope{
		Status: StatusSuccess,
		Data:   data,
	}
}

func Error(message string) Envelope {
	return Envelope{
		Status:  StatusError,
		Message: message,
	}
}