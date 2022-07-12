package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var keywords = map[string]bool{
	"struct": true,
	"func":   true,
}

func main() {

	file, err := os.Open("./exampleFile.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// limited to lines under 64k
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ParseLine(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func ParseLine(s string) {
	words := strings.Fields(s)
	for _, word := range words {
		if IsKeyword(word, keywords) {
			fmt.Println(word)
		}
	}
}

// Return whether the given word is a keyword or not
func IsKeyword(s string, k map[string]bool) bool {
	return k[s]
}
