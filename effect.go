package main

import "fmt"

type Effect struct {
	Description string
	Duration    uint64
	EndRound    uint64
}

func (e Effect) Step(round uint64) (bool, string) {
	if round >= e.EndRound {
		return true, fmt.Sprintf("EFFECT: %s", e.Description)
	}
	return false, ""
}
