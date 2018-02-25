package conf

import "strings"

func snakecaseLower(s string) string {
	return strings.ToLower(snakecase(s))
}

func snakecaseUpper(s string) string {
	return strings.ToUpper(snakecase(s))
}

func snakecase(s string) string {
	b := make([]byte, 0, 64)
	i := len(s) - 1

	// search sequences, starting from the end of the string
	for i >= 0 {
		switch {
		case isLower(s[i]): // sequence of lowercase, maybe starting with an uppercase
			for i >= 0 && !isSeparator(s[i]) && !isUpper(s[i]) {
				b = append(b, s[i])
				i--
			}

			if i >= 0 {
				b = append(b, snakebyte(s[i]))
				i--
				if isSeparator(s[i+1]) { // avoid double underscore if we have "_word"
					continue
				}
			}

			if i >= 0 && !isSeparator(s[i]) { // avoid double underscores if we have "_Word"
				b = append(b, '_')
			}

		case isUpper(s[i]): // sequence of uppercase
			for i >= 0 && !isSeparator(s[i]) && !isLower(s[i]) {
				b = append(b, s[i])
				i--
			}

			if i >= 0 {
				if isSeparator(s[i]) {
					i--
				}
				b = append(b, '_')
			}

		default: // not a letter, it'll be part of the next sequence
			b = append(b, snakebyte(s[i]))
			i--
		}
	}

	// reverse
	for i, j := 0, len(b)-1; i < j; {
		b[i], b[j] = b[j], b[i]
		i++
		j--
	}

	return string(b)
}

func snakebyte(b byte) byte {
	if isSeparator(b) {
		return '_'
	}
	return b
}

func isSeparator(c byte) bool {
	return c == '_' || c == '-'
}

func isUpper(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func isLower(c byte) bool {
	return c >= 'a' && c <= 'z'
}
