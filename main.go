package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type RoundState struct {
	Round    uint64
	Entities map[string]Entity
	Effects  map[uint64][]Effect
}

func (rs RoundState) AddEffect(desc string, duration uint64) {
	endRound := rs.Round + duration + 1
	if _, ok := rs.Effects[endRound]; !ok {
		rs.Effects[endRound] = make([]Effect, 0, 1)
	}
	effect := Effect{
		Description: desc,
		Duration:    duration,
		EndRound:    endRound,
	}
	rs.Effects[endRound] = append(rs.Effects[endRound], effect)
}

func (rs RoundState) AddEntity(name string, hp uint64) {
	entity := Entity{
		Name:        name,
		HP:          hp,
		Damage:      0,
		Afflictions: make([]Effect, 0),
	}
	rs.Entities[name] = entity
}

// bool : died
func (rs RoundState) DamageEntity(name string, damage uint64) (bool, error) {
	entity, ok := rs.Entities[name]
	if !ok {
		return false, fmt.Errorf("no entity with name %s exists", name)
	}
	entity.Damage += damage
	if entity.Damage >= entity.HP {
		delete(rs.Entities, name)
		return true, nil
	}
	return false, nil
}

func (rs RoundState) Clone() RoundState {
	newRoundState := RoundState{
		Round:    rs.Round,
		Entities: make(map[string]Entity),
		Effects:  make(map[uint64][]Effect),
	}
	for name, entity := range rs.Entities {
		newRoundState.Entities[name] = entity
	}
	for round, effects := range rs.Effects {
		newRoundState.Effects[round] = make([]Effect, 0, len(effects))
		for _, effect := range effects {
			newRoundState.Effects[round] = append(newRoundState.Effects[round], effect)
		}
	}
	return newRoundState
}

func (rs RoundState) Step(numRounds uint64) (RoundState, []string) {
	messages := make([]string, 0)
	newRoundState := RoundState{
		Round:    rs.Round + numRounds,
		Entities: make(map[string]Entity),
		Effects:  make(map[uint64][]Effect),
	}
	for name, entity := range rs.Entities {
		newEntity, msgs := entity.Step(newRoundState.Round)
		newRoundState.Entities[name] = newEntity
		messages = append(messages, msgs...)
	}
	for round, effects := range rs.Effects {
		for _, effect := range effects {
			over, msg := effect.Step(newRoundState.Round)
			if over {
				messages = append(messages, msg)
			} else {
				if _, ok := newRoundState.Effects[round]; !ok {
					newRoundState.Effects[round] = make([]Effect, 0, len(effects))
				}
				newRoundState.Effects[round] = append(newRoundState.Effects[round], effect)
			}
		}
	}
	return newRoundState, messages
}

type State struct {
	Current RoundState
	History []RoundState
}

func parseRounds(time []string) (uint64, error) {
	if len(time) != 2 {
		return 0, fmt.Errorf("invalid round time %v", time)
	}
	number, err := strconv.ParseUint(time[0], 10, 64)
	if err != nil {
		return 0, err
	}
	switch time[1] {
	case "r":
		return number, nil
	case "s":
		x := uint64(float64(number) / 10.0)
		if number == x*10 {
			return 0, fmt.Errorf("seconds %d is not divisible by 10", number)
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
		return 0, fmt.Errorf("unrecognized unit of time %s", time[1])
	}
}

func (s *State) input(text string) ([]string, error) {
	split := strings.Split(text, " ")
	if len(split) == 0 {
		return nil, fmt.Errorf("expected command, got %s", text)
	}
	messages := make([]string, 0)
	var err error
	var numRounds uint64 = 1
	switch split[0] {
	case "undo":
		s.Current = s.History[len(s.History)-1]
		s.History = s.History[:len(s.History)-1]
		messages = append(messages, "state rewinded")
		return messages, nil
	case "step":
		if len(split) > 1 {
			numRounds, err = parseRounds(split[1:])
			if err != nil {
				return nil, err
			}
		}
		newState, msgs := s.Current.Step(numRounds)
		s.History = append(s.History, s.Current)
		s.Current = newState
		return msgs, nil
	case "for":
		numRounds, err = parseRounds(split[1:])
		if err != nil {
			return nil, err
		}
		newState := s.Current.Clone()
	case "in":
	case "add":
	case "del":
	case "dmg":
	case "reset":
	default:
		return nil, fmt.Errorf("unrecognized command %s", split[0])
	}
}

func main() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println(text)
	}
}
