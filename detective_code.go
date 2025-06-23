// detective_code.go

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

/*
The data file contains one row for each cipher.
Blank rows and rows with '#' in column 1 are ignored.

Field
1       Puzzle number in decimal
2       Cipher number on page
3       Cipher type
3-n 	The cipher

The cipher types are:
A Detecive Code
B Anagrams
C Caesar
*/
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

	cipher := make(map[int]string)
	datafile, err := os.Open(DATAFILE)
	scanner := bufio.NewScanner(datafile)
	var ix, jx int
	for scanner.Scan() {
		line := scanner.Text()
		err = scanner.Err()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		lineFields := strings.Fields(line)
		ix, err = strconv.Atoi(lineFields[0])
		if err != nil {
			fmt.Println(line)
			panic("First field of line is not integer")
		}
		jx, err = strconv.Atoi(lineFields[1])
		if err != nil {
			fmt.Println(line)
			panic("Second field of line is not integer")
		}
		ix = ix*10 + jx
		// fmt.Println(ix)
		cipher[ix] = line[2:]
		cipher[ix] = strings.Join(lineFields[3:], " ")
	}
	// os.Exit(0)
	c, _ := strconv.Atoi(os.Args[1]) // get the cipher number

	// Allow the user to enter just the puzzle number without the trailing cipher number.
	if c < 100 {
		c = c*10 + 1
	}
	if len(cipher[c]) == 0 {
		panic("\nUndefined cipher.")
	}
	ciphertext := []rune(cipher[c])
	plaintext := make([]rune, len(ciphertext))
	for i, r := range ciphertext {
		r = unicode.ToUpper(r)
		if r < 'A' || r > 'Z' {
			plaintext[i] = r
		} else {
			plaintext[i] = lokup[r-'A']
		}
	}
	fmt.Println(string(plaintext))
}
