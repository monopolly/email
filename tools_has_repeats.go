package email

// ---------------------------
func HasRepeats(s string) (res bool) {
	if s == "" {
		return false
	}

	first := s[0]
	for _, x := range s {
		if x != rune(first) {
			return
		}
	}

	return true
}
