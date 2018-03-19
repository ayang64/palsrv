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

type Challenger struct {
	MinWords         int
	MaxWords         int
	MinPalindromes   int
	MaxPalindromes   int
	ChallengeLimit   int
	ChallengeCurrent int
	Progress         float64
	Words            []string
	Expecting        []string
	Answer           string
}

func (c *Challenger) HandleChallenge(conn net.Conn) {
	for c.ChallengeCurrent = 0; c.ChallengeCurrent < c.ChallengeLimit; c.ChallengeCurrent++ {
		c.SetChallenge()
		c.IssueChallenge(conn)
		log.Printf(">> %d challenges left.", c.ChallengeLimit-c.ChallengeCurrent)
		if err := c.ReceiveResponse(conn); err != nil {
			fmt.Fprintf(conn, "!!! error: %s\n", err)
			conn.Close()
			return
		}
	}
	fmt.Fprintf(conn, "!!! flag[%s]\n", c.Answer)
	conn.Close()
}

func isPalindrome(s string) bool {
	for i, r := 0, []rune(s); i < len(r)/2; i++ {
		if r[i] != r[len(r)-i-1] {
			return false
		}
	}
	return true
}

func (c *Challenger) RandKoreanPalindrome(min, max int) (rc string) {
	kchars := []rune{
		'ㄱ', '기', '역', '기', '윽', 'ㄴ', '니', '은', 'ㄷ', '디',
		'귿', '디', '읃', 'ㄹ', '리', '을', 'ㅁ', '미', '음', 'ㅂ',
		'비', '읍', 'ㅅ', '시', '옷', '시', '읏', 'ㅇ', '이', '응',
		'ㅈ', '지', '읒', 'ㅊ', '치', '읓', 'ㅋ', '키', '읔', 'ㅌ',
		'티', '읕', 'ㅍ', '피', '읖', 'ㅎ', '히', '읗',
	}

	s := rand.Intn(max-min) + min
	r := make([]rune, s)
	for i := 0; i < s/2+s%2; i++ {
		idx := rand.Intn(len(kchars))
		r[i] = rune(kchars[idx])
		r[len(r)-1-i] = r[i]
	}

	return string(r)
}

func (c *Challenger) RandPalindrome(min, max int) (rc string) {
	s := rand.Intn(max-min) + min
	r := make([]rune, s)

	for i := 0; i < s/2+s%2; i++ {
		r[i] = rune('a' + rand.Intn(26))
		r[len(r)-1-i] = r[i]
	}

	return string(r)
}

// Generate a random 'word' and avoid picking an accidental
// palindrome.
func (c *Challenger) RandWord(min, max int) (rc string) {
	for {
		rc = ""
		for i, max := 0, rand.Intn(max-min)+min; i < max; i++ {
			rc += string('a' + rand.Intn(26))
		}
		// we don't want to accidentally generate a palindrome
		if isPalindrome(rc) == false {
			break
		}
	}
	return
}

func randSpace() string {
	space := []byte{' ', '\t'}

	rc := ""

	for i, max := 0, rand.Intn(5)+1; i < max; i++ {
		flip := rand.Intn(2)
		rc += string(space[flip])
	}

	return rc
}

func (c *Challenger) IssueChallenge(w io.Writer) {
	// ramp up the probability of noisy spaces to a max of 90%
	sprob := c.Progress * .90

	output := ""
	if rand.Float64() < sprob {
		for idx := range c.Words {
			output += fmt.Sprintf("%s", c.Words[idx])

			if idx != len(c.Words) {
				output += fmt.Sprintf("%s", randSpace())
			}
		}
	} else {
		output += fmt.Sprintf("%s", strings.Join(c.Words, " "))
	}

	fmt.Printf("SENDING: %s\n", output)
	fmt.Fprintf(w, "%s\n", output)
}

func (c *Challenger) ReceiveResponse(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)

	complete := make(chan error)

	go func() {
		for len(c.Expecting) > 0 {
			if scanner.Scan() == false {
				complete <- fmt.Errorf("could not parse input.")
				return
			}

			if scanner.Text() != c.Expecting[0] {
				complete <- fmt.Errorf("bad response, expected: %s", strings.Join(c.Expecting, " "))
				return
			}
			c.Expecting = c.Expecting[1:]
		}
		complete <- nil
	}()

	return func() error {
		select {
		case err := <-complete:
			return err

		case <-time.After(250 * time.Millisecond):
			return fmt.Errorf("time out, expected: %v", strings.Join(c.Expecting, " "))
		}
	}()

}

func (c *Challenger) SetChallenge() {
	// number of words we're adding to this result set.
	// nwords := rand.Intn(c.MaxWords-c.MinWords) + c.MinWords

	// generate random words and palindromes
	am := AirMix{Min: 10, Max: 30}
	am.Init()

	p := rand.Intn(am.Min-5) + 5

	c.Words = []string(nil)
	c.Expecting = []string(nil)

	c.Progress = float64(c.ChallengeCurrent) / float64(c.ChallengeLimit)

	fmt.Printf("Progress: %f\n", c.Progress)

	for {
		v, err := am.Pick()
		if err != nil {
			break
		}

		word := func() string {
			if v < p+am.Min {
				return func() string {
					// as we near the end of the run, we want to increase the probability
					// of outputting korean to a max of 40%.
					kprob := c.Progress * .40

					if rand.Float64() < kprob {
						return c.RandKoreanPalindrome(10, 20)
					}
					return c.RandPalindrome(10, 20)
				}()
			}
			return c.RandWord(10, 20)
		}()

		c.Words = append(c.Words, word)

		if v < p+am.Min {
			c.Expecting = append(c.Expecting, word)
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":12321")
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err)
		}

		c := Challenger{ChallengeLimit: rand.Intn(5000) + 500, Answer: "A Man, A Plan, A Canal-Panama!"}
		go c.HandleChallenge(conn)
	}
}
