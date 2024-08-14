package common

type Command struct {
	Type  CommandType `json:"type"`
	Key   string      `json:"key,omitempty"`
	Value string      `json:"value,omitempty"`
	// Id is used to uniquely identify a command for sqs message deduplication
	Id string `json:"id,omitempty"`
}

type CommandType string

const (
	AddItem     CommandType = "add-item"
	DeleteItem  CommandType = "delete-item"
	GetItem     CommandType = "get-item"
	GetAllItems CommandType = "get-all-items"
)
