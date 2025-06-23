package main

import "fmt"

func Reverse(s string) string {
	n := len(s)
	runes := make([]rune, n)
	for _, rune := range s {
		n--
		runes[n] = rune
	}
	return string(runes[n:])
}

func TestReverse() {
	fmt.Println(Reverse(Reverse("Hello, 世界")))
	fmt.Println(Reverse(Reverse("The quick brown 狐 jumped over the lazy 犬")))
}
