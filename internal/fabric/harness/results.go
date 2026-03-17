package harness

type Resolution struct {
	Compatible bool     `json:"compatible"`
	Profile    Profile  `json:"profile"`
	Checked    []string `json:"checked"`
	Reason     string   `json:"reason"`
}
