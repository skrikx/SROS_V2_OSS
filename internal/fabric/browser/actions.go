package browser

func Open(url string) Session {
	return Session{URL: url, Visited: true}
}
