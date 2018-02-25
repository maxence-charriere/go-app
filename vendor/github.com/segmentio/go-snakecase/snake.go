//
// Fast snake-case implementation.
//
package snakecase

// Snakecase the given string.
func Snakecase(s string) string {
	b := make([]byte, 0, 64)
	l := len(s)
	i := 0

	// loop until we reached the end of the string
	for i < l {

		// skip leading bytes that aren't letters or numbers
		for i < l && !isWord(s[i]) {
			i++
		}

		if i < l && len(b) != 0 {
			b = append(b, '_')
		}

		// Append all leading uppercase or digits
		for i < l {
			if c := s[i]; !isHead(c) {
				break
			} else {
				b = append(b, toLower(c))
			}
			i++
		}

		// Append all trailing lowercase or digits
		for i < l {
			if c := s[i]; !isTail(c) {
				break
			} else {
				b = append(b, c)
			}
			i++
		}
	}

	return string(b)
}

func isHead(c byte) bool {
	return isUpper(c) || isDigit(c)
}

func isTail(c byte) bool {
	return isLower(c) || isDigit(c)
}

func isWord(c byte) bool {
	return isLetter(c) || isDigit(c)
}

func isLetter(c byte) bool {
	return isLower(c) || isUpper(c)
}

func isUpper(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func isLower(c byte) bool {
	return c >= 'a' && c <= 'z'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func toLower(c byte) byte {
	if isUpper(c) {
		return c + ('a' - 'A')
	}
	return c
}
