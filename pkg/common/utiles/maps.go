package utiles

// SimplifyMap, simplifies a map[string][]string to map[string]interface{}
// useful for reading HTTP headers, gRPC metadata, and other similar scenarios
// where the values are either a single string or an array of strings.
func SimplifyMap(data map[string][]string) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range data {
		if len(value) == 1 {
			result[key] = value[0]
		} else {
			result[key] = value
		}
	}
	return result
}
