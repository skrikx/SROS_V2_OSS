package web

func SearchIndex(query string, corpus []string) []string {
	out := []string{}
	for _, item := range corpus {
		if query == "" || item == query {
			out = append(out, item)
		}
	}
	return out
}
