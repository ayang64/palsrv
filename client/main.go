package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func IsPalindrome(s string) bool {
	for i, r := 0, []rune(s); i < len(r)/2; i++ {
		if r[i] != r[len(r)-i-1] {
			return false
		}
	}
	return true
}

func main() {
	server := flag.String("s", "palindromer-1469409429.us-east-1.elb.amazonaws.com:7777", "address:port of server to connect to.")
	flag.Parse()

	conn, err := net.Dial("tcp", *server)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// create a scanner that splits on lines.
	ls := bufio.NewScanner(conn)

	for ls.Scan() {

		// check for the 'answer'.
		if strings.HasPrefix(ls.Text(), "!!! flag[") {
			bra := strings.Index(ls.Text(), "[") + 1
			ket := strings.Index(ls.Text(), "]")
			fmt.Printf("%s\n", ls.Text()[bra:ket])
			break
		}

		// create a reader from the line of test we've read,
		// and use it to create a word oriented scanner.
		fmt.Fprintf(os.Stderr, "%s\n", ls.Text())

		wordreader := strings.NewReader(ls.Text())
		ws := bufio.NewScanner(wordreader)
		ws.Split(bufio.ScanWords)

		var palendromes []string

		// test each word we scan for a palindrome.  if it is a palindrome, append
		// word to our output list.
		for ws.Scan() {
			fmt.Fprintf(os.Stderr, "> %q\n", ws.Text())
			if IsPalindrome(ws.Text()) {
				palendromes = append(palendromes, ws.Text())
				fmt.Fprintf(os.Stderr, "(%q)\n", ws.Text())
			}
		}

		// join our slice of palindromes and respond to server.
		response := strings.Join(palendromes, " ")
		fmt.Fprintf(conn, "%s\n", response)
	}
}
