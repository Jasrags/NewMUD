package events

type Input struct {
	UserID string
	Text   string
}

func (i Input) Type() string {
	return "input"
}

type Message struct {
	UserID         string
	ExcludeUserIDs []string
	RoomID         string
	Text           string
}

func (m Message) Type() string {
	return "message"
}
