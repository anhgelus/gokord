package utils

// AStringContains checks if a slice "a" contains a "b" string
func AStringContains(a []string, b string) bool {
	for _, c := range a {
		if c == b {
			return true
		}
	}
	return false
}
