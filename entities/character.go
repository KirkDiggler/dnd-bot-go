package entities

import "github.com/KirkDiggler/dnd-bot-go/repositories/character"

type Character struct {
	ID    string `json:"id"`
	Name  string
	Race  *Race
	Class *Class
}

func (c *Character) ToData() *character.Data {
	var raceKey string
	var classKey string

	if c.Race != nil {
		raceKey = c.Race.Key
	}

	if c.Class != nil {
		classKey = c.Class.Key
	}

	return &character.Data{
		ID:       c.ID,
		Name:     c.Name,
		RaceKey:  raceKey,
		ClassKey: classKey,
	}
}
