package main

import (
	"math/rand"
	"os"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		Name     string
		Value    string
		Expected bool
	}{
		{Name: "ASCII Palindrome", Value: "abcdcba", Expected: true},
		{Name: "ASCII Non-Palindrome", Value: "xabcdcba", Expected: false},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result := IsPalindrome(test.Value)
			if result != test.Expected {
				t.Fatalf("IsPalindrome(%q) returned %v; expected %v", test.Value, result, test.Expected)
			}
		})
	}
}

func TestChallengerSetChallenge(t *testing.T) {

	c := Challenger{}

	c.SetChallenge()
	c.IssueChallenge(os.Stdout)
}

func TestAirMix(t *testing.T) {
	am := AirMix{Min: 5, Max: 10}
	am.Init()

	t.Logf("%#v", am.Ball)

	for {
		v, err := am.Pick()

		if err != nil {
			break
		}

		t.Logf("%d", v)
	}
}

func TestRandomPalindrome(t *testing.T) {
	for i := 0; i < 100; i++ {
		rc := RandPalindrome(7, 30)

		if IsPalindrome(rc) == false {
			t.Fatalf("IsPalindrome(%q) returned %v; expected true!", rc, false)
		}
		t.Logf("%s", rc)
	}
}

func TestRandomWord(t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Logf("%s", RandomWord(7, 30))
	}
}
