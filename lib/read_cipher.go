package lib

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
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

func ReadCipher(filename string, key int) (string, error) {
	// Allow the user to enter just the puzzle number without the trailing cipher number.
	if key < 100 {
		key = key*10 + 1
	}
	datafile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
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
		if ix == key {

			return strings.Join(lineFields[3:], " "), nil
		}
	}
	return "", fmt.Errorf("No find cipher %d", key)

}
