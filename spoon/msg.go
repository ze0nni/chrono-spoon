package spoon

type PushMsg struct {
	Command string `json:"command"`

	Name  string `json:"name"`
	Group string `json:"group"`
	Start int64  `json:"start"`
	End   int64  `json:"end"`
}
