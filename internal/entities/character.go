package entities

import (
	"fmt"
	"strings"
	"sync"

	"github.com/KirkDiggler/dnd-bot-go/internal/dice"
)

type Slot string

const (
	SlotMainHand  Slot = "main-hand"
	SlotOffHand   Slot = "off-hand"
	SlotTwoHanded Slot = "two-handed"
	SlotBody      Slot = "body"
	SlotNone      Slot = "none"
)

type Character struct {
	ID                 string
	OwnerID            string
	Name               string
	Speed              int
	Race               *Race
	Class              *Class
	Attribues          map[Attribute]*AbilityScore
	Rolls              []*dice.RollResult
	Proficiencies      map[ProficiencyType][]*Proficiency
	ProficiencyChoices []*Choice
	Inventory          map[string][]Equipment

	HitDie           int
	AC               int
	MaxHitPoints     int
	CurrentHitPoints int
	Level            int
	Experience       int
	NextLevel        int

	EquippedSlots map[Slot]Equipment

	mu sync.Mutex
}

func (c *Character) Equip(e Equipment) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.EquippedSlots == nil {
		c.EquippedSlots = make(map[Slot]Equipment)
	}

	if e.GetSlot() == SlotTwoHanded {
		if c.EquippedSlots[SlotMainHand] != nil {
			c.Unequip(c.EquippedSlots[SlotMainHand])
		}
		if c.EquippedSlots[SlotOffHand] != nil {
			c.Unequip(c.EquippedSlots[SlotOffHand])
		}
	}

	// if we are trying to equip another main hand we will assign it to the off hand
	if e.GetSlot() == SlotMainHand {
		if c.EquippedSlots[SlotMainHand] != nil {
			c.EquippedSlots[SlotOffHand] = c.EquippedSlots[SlotMainHand]
		}
	}

	c.EquippedSlots[e.GetSlot()] = e
	c.calculateAC()
}

func (c *Character) calculateAC() {
	c.AC = 10
	for _, e := range c.EquippedSlots {
		if e == nil {
			continue
		}

		if e.GetEquipmentType() == "Armor" {
			armor := e.(*Armor)
			if armor.ArmorClass == nil {
				continue
			}
			if e.GetSlot() == SlotBody {
				c.AC = armor.ArmorClass.Base
				if armor.ArmorClass.DexBonus {
					// TODO: load max and bonus and limit id applicable
					c.AC += c.Attribues[AttributeDexterity].Bonus
				}
				continue
			}

			c.AC += armor.ArmorClass.Base
			if armor.ArmorClass.DexBonus {
				c.AC += c.Attribues[AttributeDexterity].Bonus
			}
		}
	}
}

func (c *Character) Unequip(e Equipment) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.EquippedSlots == nil {
		c.EquippedSlots = make(map[Slot]Equipment)
	}

	c.EquippedSlots[e.GetSlot()] = nil
}

func (c *Character) SetHitpoints() {
	if c.Attribues == nil {
		return
	}

	if c.Attribues[AttributeConstitution] == nil {
		return
	}

	if c.HitDie == 0 {
		return
	}

	c.MaxHitPoints = c.HitDie + c.Attribues[AttributeConstitution].Bonus
	c.CurrentHitPoints = c.MaxHitPoints
}
func (c *Character) AddAttribute(attr Attribute, score int) {
	if c.Attribues == nil {
		c.Attribues = make(map[Attribute]*AbilityScore)
	}

	bonus := 0
	if _, ok := c.Attribues[attr]; ok {
		bonus = c.Attribues[attr].Bonus
	}
	abilityScore := &AbilityScore{
		Score: score,
		Bonus: bonus,
	}
	switch {
	case score == 1:
		abilityScore.Bonus += -5
	case score < 4 && score > 1:
		abilityScore.Bonus += -4
	case score < 6 && score > 3:
		abilityScore.Bonus += -3
	case score < 8 && score > 5:
		abilityScore.Bonus += -2
	case score < 10 && score >= 8:
		abilityScore.Bonus += -1
	case score < 12 && score > 9:
		abilityScore.Bonus += 0
	case score < 14 && score > 11:
		abilityScore.Bonus += 1
	case score < 16 && score > 13:
		abilityScore.Bonus += 2
	case score < 18 && score > 15:
		abilityScore.Bonus += 3
	case score < 20 && score > 17:
		abilityScore.Bonus += 4
	case score == 20:
		abilityScore.Bonus += 5
	}

	c.Attribues[attr] = abilityScore
}
func (c *Character) AddAbilityBonus(ab *AbilityBonus) {
	if c.Attribues == nil {
		c.Attribues = make(map[Attribute]*AbilityScore)
	}

	if _, ok := c.Attribues[ab.Attribute]; !ok {
		c.Attribues[ab.Attribute] = &AbilityScore{}
	}

	c.Attribues[ab.Attribute] = c.Attribues[ab.Attribute].AddBonus(ab.Bonus)
}

func (c *Character) AddInventory(e Equipment) {
	if c.Inventory == nil {
		c.Inventory = make(map[string][]Equipment)
	}

	c.mu.Lock()
	if c.Inventory[e.GetEquipmentType()] == nil {
		c.Inventory[e.GetEquipmentType()] = make([]Equipment, 0)
	}

	c.Inventory[e.GetEquipmentType()] = append(c.Inventory[e.GetEquipmentType()], e)
	c.mu.Unlock()
}

func (c *Character) AddProficiency(p *Proficiency) {
	if c.Proficiencies == nil {
		c.Proficiencies = make(map[ProficiencyType][]*Proficiency)
	}
	c.mu.Lock()
	if c.Proficiencies[p.Type] == nil {
		c.Proficiencies[p.Type] = make([]*Proficiency, 0)
	}

	c.Proficiencies[p.Type] = append(c.Proficiencies[p.Type], p)
	c.mu.Unlock()
}

func (c *Character) AddAbilityScoreBonus(attr Attribute, bonus int) {
	if c.Attribues == nil {
		c.Attribues = make(map[Attribute]*AbilityScore)
	}

	c.Attribues[attr] = c.Attribues[attr].AddBonus(bonus)
}

func (c *Character) String() string {
	msg := strings.Builder{}
	if c.Race == nil || c.Class == nil {
		return "Character not fully created"
	}

	msg.WriteString(fmt.Sprintf("%s the %s %s\n", c.Name, c.Race.Name, c.Class.Name))

	msg.WriteString("**Rolls**:\n")
	for _, roll := range c.Rolls {
		msg.WriteString(fmt.Sprintf("%s, ", roll))
	}
	msg.WriteString("\n")
	msg.WriteString("\n**Stats**:\n")
	msg.WriteString(fmt.Sprintf("  -  Speed: %d\n", c.Speed))
	msg.WriteString(fmt.Sprintf("  -  Hit Die: %d\n", c.HitDie))
	msg.WriteString(fmt.Sprintf("  -  AC: %d\n", c.AC))
	msg.WriteString(fmt.Sprintf("  -  Max Hit Points: %d\n", c.MaxHitPoints))
	msg.WriteString(fmt.Sprintf("  -  Current Hit Points: %d\n", c.CurrentHitPoints))
	msg.WriteString(fmt.Sprintf("  -  Level: %d\n", c.Level))
	msg.WriteString(fmt.Sprintf("  -  Experience: %d\n", c.Experience))

	msg.WriteString("\n**Attributes**:\n")
	for _, attr := range Attributes {
		if c.Attribues[attr] == nil {
			continue
		}
		msg.WriteString(fmt.Sprintf("  -  %s: %s\n", attr, c.Attribues[attr]))
	}

	msg.WriteString("\n**Proficiencies**:\n")
	for _, key := range ProficiencyTypes {
		if c.Proficiencies[key] == nil {
			continue
		}

		msg.WriteString(fmt.Sprintf("  -  **%s**:\n", key))
		for _, prof := range c.Proficiencies[key] {
			msg.WriteString(fmt.Sprintf("    -  %s\n", prof.Name))
		}
	}

	msg.WriteString("\n**Inventory**:\n")
	for key := range c.Inventory {
		if c.Inventory[key] == nil {
			continue
		}

		msg.WriteString(fmt.Sprintf("  -  **%s**:\n", key))
		for _, item := range c.Inventory[key] {
			msg.WriteString(fmt.Sprintf("    -  %s\n", item.GetName()))
		}

	}
	return msg.String()
}
