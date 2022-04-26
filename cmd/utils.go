package cmd

// isPresent will verify whether a string is present in a slice.
// Returns true if yes, false otherwise.
func isPresent(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
