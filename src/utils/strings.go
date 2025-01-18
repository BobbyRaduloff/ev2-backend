package utils

func FilterEmptyStrings(strings []string) []string {
	var result []string
	for _, str := range strings {
		if str != "" {
			result = append(result, str)
		}
	}
	return result
}
