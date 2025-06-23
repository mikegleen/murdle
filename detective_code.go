// detective_code.go

package main

import (
	"fmt"
	"lib"
	"os"
	"strconv"
	"unicode"
)

const DATAFILE = "data.txt"

func main() {

	if len(os.Args) < 2 {
		panic("\nOne parameter required, the cipher number.")
	}

	// Translate table for the Murdle code "A" where the cipher alphabet is in reverse order.
	lokup := make([]rune, 26)
	for r := 'A'; r <= 'Z'; r++ {
		ix := 'Z' - r
		lokup[ix] = r
	}
	// fmt.Println(lokup)

	c, _ := strconv.Atoi(os.Args[1]) // get the cipher number

	ciphertext, err := lib.ReadCipher(DATAFILE, c)
	if err != nil {
		panic(err)
	}
	ciphertextr := []rune(ciphertext)
	plaintext := make([]rune, len(ciphertextr))
	for i, r := range ciphertextr {
		r = unicode.ToUpper(r)
		if r < 'A' || r > 'Z' {
			plaintext[i] = r
		} else {
			plaintext[i] = lokup[r-'A']
		}
	}
	fmt.Println(string(plaintext))
}
