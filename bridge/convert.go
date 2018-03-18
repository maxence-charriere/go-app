package bridge

// Strings is an helper function that converts the given value to a slice of
// strings.
// It panics if v is not a slice of interface{}.
func Strings(v interface{}) []string {
	src := v.([]interface{})
	s := make([]string, 0, len(src))

	for _, item := range src {
		s = append(s, item.(string))
	}
	return s
}
