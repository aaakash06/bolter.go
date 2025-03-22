package utils

import "strings"

// stripIndents takes a string or a slice of strings and returns a single string with indents stripped.
func stripIndents(arg0 interface{}, values ...interface{}) string {
	switch v := arg0.(type) {
	case string:
		return _stripIndents(v)
	case []string:
		var processedString strings.Builder
		for i, curr := range v {
			processedString.WriteString(curr)
			if i < len(values) {
				processedString.WriteString(values[i].(string))
			}
		}
		return _stripIndents(processedString.String())
	default:
		return ""
	}
}

// _stripIndents removes leading and trailing whitespace from each line of the input string.
func _stripIndents(value string) string {
	lines := strings.Split(value, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	return strings.Join(lines, "\n")
}
