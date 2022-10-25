package tracker

import (
	"fmt"
	"sort"
	"strings"
)

type Entity struct {
	Name        string
	HP          uint64
	Damage      uint64
	Afflictions []Effect
}

func (e Entity) Step(round uint64) (Entity, []string) {
	messages := make([]string, 0)
	newEntity := Entity{
		Name:        e.Name,
		HP:          e.HP,
		Damage:      e.Damage,
		Afflictions: make([]Effect, 0, len(e.Afflictions)),
	}
	for _, affliction := range e.Afflictions {
		over, msg := affliction.Step(round)
		if over {
			messages = append(messages, msg)
		} else {
			newEntity.Afflictions = append(newEntity.Afflictions, affliction)
		}
	}
	return newEntity, messages
}

func (e Entity) Normalize(by uint64) Entity {
	normalizedAfflictions := make([]Effect, 0, len(e.Afflictions))
	for _, effect := range e.Afflictions {
		normalizedAfflictions = append(normalizedAfflictions, effect.Normalize(by))
	}
	return Entity{
		Name:        e.Name,
		HP:          e.HP,
		Damage:      e.Damage,
		Afflictions: normalizedAfflictions,
	}
}

func (e Entity) Afflict(until uint64, duration uint64, desc string) Entity {
	e.Afflictions = append(e.Afflictions, Effect{
		Description: fmt.Sprintf("%s (%s)", desc, e.Name),
		Duration:    duration,
		EndRound:    until,
	})
	return e
}

func (e Entity) String() string {
	messages := []string{fmt.Sprintf("%s | HP:%d, Damage:%d, Afflictions:", e.Name, e.HP, e.Damage)}
	afflictionMessages := make([]string, 0, len(e.Afflictions))
	for _, affliction := range e.Afflictions {
		afflictionMessages = append(afflictionMessages, fmt.Sprintf("  %s", affliction.String()))
	}
	sort.Strings(afflictionMessages)
	messages = append(messages, afflictionMessages...)
	return strings.Join(messages, "\n")
}
