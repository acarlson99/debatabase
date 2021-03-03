package main

func filterArr(arr []string, f func(string) bool) []string {
	r := []string{}
	for _, s := range arr {
		if f(s) {
			r = append(r, s)
		}
	}
	return r
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
