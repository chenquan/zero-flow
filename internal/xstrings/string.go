package xstrings

func Distinct(values []string) []string {
	m := make(map[string]struct{}, len(values))
	newValues := make([]string, 0, len(values))
	for _, value := range values {
		_, ok := m[value]
		if ok {
			continue
		}
		newValues = append(newValues, value)
		m[value] = struct{}{}
	}

	return newValues
}
