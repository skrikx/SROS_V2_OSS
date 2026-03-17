package browser

type Session struct {
	URL     string `json:"url"`
	Visited bool   `json:"visited"`
}
