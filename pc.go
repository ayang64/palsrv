package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

type AirMix struct {
	Min  int
	Max  int
	Ball []int
}

func (a *AirMix) Init() error {
	if a.Min >= a.Max {
		return fmt.Errorf("min must be greater than max")
	}

	// create balls
	a.Ball = make([]int, 0, a.Max-a.Min)

	for i := a.Min; i < a.Max; i++ {
		a.Ball = append(a.Ball, i)
	}
	return nil
}

func (a *AirMix) Pick() (int, error) {
	if len(a.Ball) == 0 {
		return 0, fmt.Errorf("out of balls")
	}

	i := rand.Intn(len(a.Ball))
	rc := a.Ball[i]
	a.Ball = append(a.Ball[0:i], a.Ball[i+1:]...)

	return rc, nil
}

type Challenger struct {
	MinWords       int
	MaxWords       int
	MinPalindromes int
	MaxPalindromes int
	ChallengeLimit int
	Words          []string
	Expecting      []string
	Answer         string
}

func (c *Challenger) HandleChallenge(conn net.Conn) {
	for i := 0; i < c.ChallengeLimit; i++ {
		c.SetChallenge()
		c.IssueChallenge(conn)
		log.Printf(">> %d challenges left.", c.ChallengeLimit-i)
		log.Printf(">> %v\n", c.Expecting)
		if err := c.ReceiveResponse(conn); err != nil {
			fmt.Fprintf(conn, "!!! ERROR: %s\n", err)
			conn.Close()
			return
		}
	}
	fmt.Fprintf(conn, "!!! flag[%s]\n", c.Answer)
	conn.Close()
}

func IsPalindrome(s string) bool {
	for i, r := 0, []rune(s); i < len(r)/2; i++ {
		if r[i] != r[len(r)-i-1] {
			return false
		}
	}
	return true
}

func RandPalindrome(min, max int) (rc string) {
	s := rand.Intn(max-min) + min
	r := make([]rune, s)

	for i := 0; i < s/2+s%2; i++ {
		r[i] = rune('a' + rand.Intn(26))
		r[len(r)-1-i] = r[i]
	}

	return string(r)
}

func RandomWord(min, max int) (rc string) {
	for {
		rc = ""
		for i, max := 0, rand.Intn(max-min)+min; i < max; i++ {
			rc += string('a' + rand.Intn(26))
		}
		// we don't want to accidentally generate a palindrome
		if IsPalindrome(rc) == false {
			break
		}
	}
	return
}

func (c *Challenger) IssueChallenge(w io.Writer) {
	fmt.Fprintf(w, "%s\n", strings.Join(c.Words, " "))
}

func (c *Challenger) ReceiveResponse(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)

	complete := make(chan error)

	go func() {
		for len(c.Expecting) > 0 {
			if scanner.Scan() == false {
				break
			}

			if scanner.Text() != c.Expecting[0] {
				log.Printf("bad response!\n")
				complete <- fmt.Errorf("EXPECTED %s", strings.Join(c.Expecting, " "))
			}

			log.Printf("good response!\n")
			c.Expecting = c.Expecting[1:]
		}
		complete <- nil
	}()

	select {
	case err := <-complete:
		return err

	case <-time.After(500 * time.Millisecond):
		return fmt.Errorf("TIME OUT")
	}

	return nil
}

func (c *Challenger) SetChallenge() {
	// number of words we're adding to this result set.
	// nwords := rand.Intn(c.MaxWords-c.MinWords) + c.MinWords

	// generate random words and palindromes

	am := AirMix{Min: 3, Max: 30}
	am.Init()

	p := rand.Intn(am.Min-2) + 2

	c.Words = []string(nil)
	c.Expecting = []string(nil)

	for {
		v, err := am.Pick()
		if err != nil {
			break
		}

		word := func() string {
			if v < p+am.Min {
				return RandPalindrome(7, 30)
			}
			return RandomWord(7, 30)
		}()

		c.Words = append(c.Words, word)

		if v < p+am.Min {
			c.Expecting = append(c.Expecting, word)
		}
	}
}

func main() {
	rand.Seed(time.Now().Unix())

	ln, err := net.Listen("tcp", ":12321")
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err)

		}
		c := Challenger{ChallengeLimit: rand.Intn(10000) + 10000, Answer: "ORLANDO GOPHERS IS FUN!"}
		go c.HandleChallenge(conn)
	}
}
