package util

import "fmt"

// Parse xxxyyyzzz to xxx.yyy.zzz for version numbers
func VerFromDec(version int) string {
	major, remainder := Divmod(version, 1_000_000)
	minor, remainder := Divmod(remainder, 1000)
	patch := remainder % 1000
	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}

func Divmod(a, b int) (int, int) {
	return a / b, a % b
}
