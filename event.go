package thunderbird

type Event struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Body    string `json:"body"`
}
