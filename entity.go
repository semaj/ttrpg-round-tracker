package main

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
		Type:        e.Type,
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
