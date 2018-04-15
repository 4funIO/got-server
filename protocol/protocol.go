package protocol

const (
	ClientVersionMin = 1097
	ClientVersionMax = 1098
	ClientVersionStr = "10.98"
)

type Protocol interface {
	ReceiveMessage(msg *NetworkMessage) error
}

type ErrDisconnectUser struct {
	Message string
	Version uint16
}

func (e ErrDisconnectUser) Error() string {
	return e.Message
}
