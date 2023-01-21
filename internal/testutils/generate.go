package testutils

func GenerateStringLongerThanMaxLength(maxLength int) string {
	var s string
	for i := 0; i < maxLength+1; i++ {
		s += "s"
	}
	return s
}
