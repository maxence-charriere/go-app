package core

import "fmt"

// ConvertToStringSlice convert the given value to string slice. It returns nil
// when the conversion is not possible.
func ConvertToStringSlice(v interface{}) []string {
	if slice, ok := v.([]string); ok {
		return slice
	}

	slice, ok := v.([]interface{})
	if !ok {
		return nil
	}

	strings := make([]string, 0, len(slice))

	for _, s := range slice {
		switch v := s.(type) {
		case string:
			strings = append(strings, v)

		default:
			strings = append(strings, fmt.Sprint(v))
		}
	}

	return strings
}
