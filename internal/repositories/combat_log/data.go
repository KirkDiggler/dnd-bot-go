package combat_log

import (
	"github.com/KirkDiggler/dnd-bot-go/internal/entities/attack"
	"time"
)

type Data struct {
	ID          string         `json:"id"`
	Type        Type           `json:"type"`
	EncounterID string         `json:"encounter_id"`
	PlayerID    string         `json:"player_id"`
	MonsterID   string         `json:"monster_id"`
	RoomID      string         `json:"room_id"`
	CreatedAt   time.Time      `json:"created_at"`
	AttackRoll  *attack.Result `json:"attack_roll"`
}

type Type string

const (
	TypeUnset         Type = ""
	TypePlayerAttack  Type = "player_attack"
	TypeMonsterAttack Type = "monster_attack"
)
