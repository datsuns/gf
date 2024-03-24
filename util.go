package main

func TrimLastOne(s string) string {
	if len(s) == 0 {
		return ""
	} else {
		return s[:len(s)-1]
	}
}
