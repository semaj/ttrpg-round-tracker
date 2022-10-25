package tracker

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type State struct {
	Current RoundState
	History []RoundState
}

func NewState() *State {
	return &State{
		Current: NewRoundState(),
		History: make([]RoundState, 0),
	}
}

var regex = regexp.MustCompile(`(\d+)(\w)`)

func parseRounds(time string) (uint64, error) {
	submatches := regex.FindSubmatch([]byte(time))
	if len(submatches) != 3 {
		return 0, fmt.Errorf("invalid rounds time %s", time)
	}
	rawNumber := string(submatches[1])
	duration := string(submatches[2])
	number, err := strconv.ParseUint(rawNumber, 10, 64)
	if err != nil {
		return 0, err
	}
	switch duration {
	case "r":
		return number, nil
	case "s":
		x := uint64(float64(number) / 10.0)
		if number != x*10 {
			return 0, fmt.Errorf("seconds %d is not divisible by 10 (%d)", number, x)
		}
		return x, nil
	case "m":
		return number * 6, nil
	case "h":
		return number * 6 * 60, nil
	case "d":
		return number * 6 * 60 * 24, nil
	case "y":
		return number * 6 * 60 * 365, nil
	default:
		return 0, fmt.Errorf("unrecognized unit of time %s", duration)
	}
}

func (s *State) String() string {
	return s.Current.String()
}

func (s *State) CurrentRound() uint64 {
	return s.Current.Round
}

func (s *State) Input(text string) ([]string, error) {
	split := strings.Split(text, " ")
	if len(split) == 0 {
		return nil, fmt.Errorf("expected command, got %s", text)
	}
	messages := make([]string, 0)
	var err error
	var numRounds uint64 = 1
	switch split[0] {
	case "print":
		messages = append(messages, s.String())
		return messages, nil
	case "undo":
		s.Current = s.History[len(s.History)-1]
		s.History = s.History[:len(s.History)-1]
		messages = append(messages, "state rewinded")
		return messages, nil
	case "step":
		if len(split) > 1 {
			numRounds, err = parseRounds(split[1])
			if err != nil {
				return nil, err
			}
		}
		newState, msgs := s.Current.Step(numRounds)
		s.History = append(s.History, s.Current)
		s.Current = newState
		return msgs, nil
	case "in":
		if len(split) < 3 {
			return nil, fmt.Errorf("command `in` requires 2 arguments, got %d", len(split)-1)
		}
		numRounds, err = parseRounds(split[1])
		if err != nil {
			return nil, err
		}
		newState := s.Current.Clone()
		newState.AddEffect(numRounds, strings.Join(split[2:], " "))
		s.History = append(s.History, s.Current)
		s.Current = newState
	case "add":
		if len(split) < 3 {
			return nil, fmt.Errorf("command `add` requires 2 arguments, got %d", len(split)-1)
		}
		hp, err := strconv.ParseUint(split[2], 10, 64)
		if err != nil {
			return nil, err
		}
		newState := s.Current.Clone()
		if err := newState.AddEntity(split[1], hp); err != nil {
			return nil, err
		}
		s.History = append(s.History, s.Current)
		s.Current = newState
	case "normalize":
		newState := s.Current.Normalize()
		s.History = append(s.History, s.Current)
		s.Current = newState
	case "afflict":
		if len(split) < 4 {
			return nil, fmt.Errorf(
				"command `afflict` requires at least 3 arguments, got %d",
				len(split)-1,
			)
		}
		numRounds, err = parseRounds(split[2])
		if err != nil {
			return nil, err
		}
		newState := s.Current.Clone()
		if err := newState.AfflictEntity(split[1], numRounds, strings.Join(split[3:], " ")); err != nil {
			return nil, err
		}
		s.History = append(s.History, s.Current)
		s.Current = newState
	case "del":
		if len(split) < 2 {
			return nil, fmt.Errorf("command `del` requires 1 arguments, got %d", len(split)-1)
		}
		newState := s.Current.Clone()
		if err := newState.DeleteEntity(split[1]); err != nil {
			return nil, err
		}
		s.History = append(s.History, s.Current)
		s.Current = newState
	case "dmg":
		if len(split) < 3 {
			return nil, fmt.Errorf("command `dmg` requires 2 arguments, got %d", len(split)-1)
		}
		damage, err := strconv.ParseUint(split[2], 10, 64)
		if err != nil {
			return nil, err
		}
		newState := s.Current.Clone()
		dead, err := newState.DamageEntity(split[1], damage)
		if err != nil {
			return nil, err
		}
		if dead {
			messages = append(messages, fmt.Sprintf("entity %s is dead!", split[1]))
		}
		s.History = append(s.History, s.Current)
		s.Current = newState
	default:
		return nil, fmt.Errorf("unrecognized command %s", split[0])
	}
	return messages, nil
}
