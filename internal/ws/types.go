package ws

type MessageData struct {
	MessageId   	int     `json:"messageId" mapstructure:"messageId"`
	Owner 			string  `json:"owner" mapstructure:"owner"`
	ConnectionId 	string 	`json:"connectionId" mapstructure:"connectionId"`
	Content 		string	`json:"content" mapstructure:"content"`
	Status 			int 	`json:"status"  mapstructure:"status"`
	CreatedTime 	int64   `json:"createdTime" mapstructure:"createdTime"`
}

type Payload struct {
	// username that sends payload
	UserName        string              `json:"username" mapstructure:"username"`
	ConnectionID 	uint32 				`json:"connectionId" mapstructure:"connectionId"`
	MessageType 	int 				`json:"messageType" mapstructure:"messageType"`
	Data 			*MessageData 		`json:"data" mapstructure:"data"`
}

const (
	SendMessage = 0
	ReceiveMessage = 1
	DeleteMessage = 2
	NewConnection = 3
	ErrorMessage = 4
)