package tracker

import "fmt"

type Effect struct {
	Description string
	Duration    uint64
	EndRound    uint64
}

func (e Effect) Normalize(by uint64) Effect {
	return Effect{
		Description: e.Description,
		Duration:    e.Duration,
		EndRound:    e.EndRound - by,
	}
}

func (e Effect) Step(round uint64) (bool, string) {
	if round >= e.EndRound {
		return true, fmt.Sprintf("EFFECT: %s", e.Description)
	}
	return false, ""
}

func (e Effect) String() string {
	return fmt.Sprintf("at %d: %s", e.EndRound, e.Description)
}
