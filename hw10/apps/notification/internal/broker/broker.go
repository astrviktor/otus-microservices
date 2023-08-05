package broker

type InterfaceBroker interface {
	Connect() error
	Send(message *Message) error
	Receive() (*Message, error)
	Close() error
}

type Message struct {
	ClientID int64  `json:"client_id"`
	OrderID  int64  `json:"order_id"`
	Theme    string `json:"theme"`
	Message  string `json:"message"`
}
