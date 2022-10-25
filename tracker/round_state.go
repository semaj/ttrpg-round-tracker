package tracker

import (
	"fmt"
	"sort"
	"strings"
)

type RoundState struct {
	Round    uint64
	Entities map[string]Entity
	Effects  map[uint64][]Effect
}

func NewRoundState() RoundState {
	return RoundState{
		Round:    1,
		Entities: make(map[string]Entity),
		Effects:  make(map[uint64][]Effect),
	}
}

func (rs RoundState) Normalize() RoundState {
	normalizedEntities := make(map[string]Entity)
	for name, entity := range rs.Entities {
		normalizedEntities[name] = entity.Normalize(rs.Round - 1)
	}
	normalizedEffects := make(map[uint64][]Effect)
	for round, effects := range rs.Effects {
		normalizedEffects[round] = make([]Effect, 0, len(effects))
		for _, effect := range effects {
			normalizedEffects[round] = append(
				normalizedEffects[round],
				effect.Normalize(rs.Round-1),
			)
		}
	}
	return RoundState{
		Round:    1,
		Entities: normalizedEntities,
		Effects:  normalizedEffects,
	}
}

func (rs RoundState) String() string {
	messages := []string{fmt.Sprintf("-Round: %d-", rs.Round)}
	messages = append(messages, "Entities:")
	entityMessages := make([]string, 0, len(rs.Entities))
	for _, entity := range rs.Entities {
		entityMessages = append(entityMessages, fmt.Sprintf(" %s", entity.String()))
	}
	sort.Strings(entityMessages)
	messages = append(messages, entityMessages...)
	messages = append(messages, "Effects:")
	effectMessages := make([]string, 0, len(rs.Effects))
	for _, effects := range rs.Effects {
		for _, effect := range effects {
			effectMessages = append(effectMessages, fmt.Sprintf(" %s", effect.String()))
		}
	}
	sort.Strings(effectMessages)
	messages = append(messages, effectMessages...)
	return strings.Join(messages, "\n")
}

func (rs RoundState) AddEffect(duration uint64, desc string) {
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

func (rs RoundState) AddEntity(name string, hp uint64) error {
	if _, ok := rs.Entities[name]; ok {
		return fmt.Errorf("entity with name %s already exists", name)
	}
	entity := Entity{
		Name:        name,
		HP:          hp,
		Damage:      0,
		Afflictions: make([]Effect, 0),
	}
	rs.Entities[name] = entity
	return nil
}

func (rs RoundState) DeleteEntity(name string) error {
	if _, ok := rs.Entities[name]; !ok {
		return fmt.Errorf("no entity with name %s exists", name)
	}
	delete(rs.Entities, name)
	return nil
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
	rs.Entities[name] = entity
	return false, nil
}

// bool : died
func (rs RoundState) AfflictEntity(name string, numRounds uint64, desc string) error {
	entity, ok := rs.Entities[name]
	if !ok {
		return fmt.Errorf("no entity with name %s exists", name)
	}
	endRound := rs.Round + numRounds + 1
	rs.Entities[name] = entity.Afflict(endRound, numRounds, desc)
	return nil
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
		newRoundState.Effects[round] = append(newRoundState.Effects[round], effects...)
	}
	return newRoundState
}

func (rs RoundState) Step(numRounds uint64) (RoundState, []string) {
	newRoundState := RoundState{
		Round:    rs.Round + numRounds,
		Entities: make(map[string]Entity),
		Effects:  make(map[uint64][]Effect),
	}
	messages := make([]string, 0)
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
