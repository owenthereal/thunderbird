package thunderbird

type Event struct {
	Command string `json:"command"`
	Channel string `json:"channel"`
	Body    string `json:"body"`
}
