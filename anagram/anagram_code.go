// detective_code.go

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/etnz/permute"
)

/*
The data file contains one row for each cipher.
Blank rows and rows with '#' in column 1 are ignored.

Columns
1-2     Puzzle number in decimal
3-n 	The cipher
*/
const DATAFILE = "../data.txt"
const DICTFILE = "/Users/mlg/pyprj/caesar/data/dictionary.txt"

// func nextline(scanner bufio.Scanner) (string, error) {

// }

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
	if err != nil {
		panic(fmt.Sprint("Cannot open ", DATAFILE))
	}
	scanner := bufio.NewScanner(datafile)
	var ix int
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
		fmt.Sscanf(line, "%2d", &ix)
		// fmt.Println(ix)
		_, ok := cipher[ix]
		if ok {
			panic(fmt.Sprint("Duplicate problem number: ", ix))
		}
		// fmt.Println(ix)
		cipher[ix] = line[2:]
	}
	datafile.Close()

	/*
		Load the dictionary
	*/
	wordDict := make(map[string]struct{})
	datafile, err = os.Open(DICTFILE)
	if err != nil {
		panic(fmt.Sprint("Cannot open ", DATAFILE))
	}
	scanner = bufio.NewScanner(datafile)
	for scanner.Scan() {
		line := scanner.Text()
		err = scanner.Err()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading dictionary.")
			log.Fatal(err)
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		wordDict[line] = struct{}{}
	}
	datafile.Close()

	// os.Exit(0)
	c, _ := strconv.Atoi(os.Args[1]) // get the cipher number
	if len(cipher[c]) == 0 {
		panic("\nUndefined cipher.")
	}

	reg, _ := regexp.Compile("[^A-Z]+") // remove everything except letters
	words := strings.Fields(cipher[c])
	// fmt.Println("cipher:", cipher[c], "\nwords: ", words)
	for _, word := range words {
		rword := reg.ReplaceAllString(word, "")
		fmt.Println("word: ", rword)
		guesses := make(map[string]struct{})
		w := []rune(rword)
		for _, try := range permute.Permutations(w) {
			stry := string(try)
			// If there are repeated letters, the same guess will be printed multiple times.
			_, found := guesses[stry]
			if found {
				continue
			}
			// fmt.Println("try: ", stry)
			_, found = wordDict[stry]
			if found {
				guesses[stry] = struct{}{}
				fmt.Printf("%16s %s\n", "", stry)
			}
		}
		// fmt.Println(guesses)
		if len(guesses) == 0 {
			fmt.Printf("%16s %s\n", "", "?")
		}
	}
}
