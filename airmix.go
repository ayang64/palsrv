package main

import (
	"fmt"
	"math/rand"
	"time"
)

type AirMix struct {
	rand *rand.Rand
	Min  int
	Max  int
	Ball []int
}

func (a *AirMix) Init() error {
	a.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
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

	i := a.rand.Intn(len(a.Ball))
	rc := a.Ball[i]
	a.Ball = append(a.Ball[0:i], a.Ball[i+1:]...)

	return rc, nil
}
