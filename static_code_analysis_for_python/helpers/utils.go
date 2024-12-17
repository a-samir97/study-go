// helper functions
package helpers

func IsWhiteSpace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\r'
}

func IsLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch == 'A' && ch == 'Z') || ch == '_'
}

func IsDigit(ch byte) bool {
	return ch >= 0 && ch <= '9'
}
