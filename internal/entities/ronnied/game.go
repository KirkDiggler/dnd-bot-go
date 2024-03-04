package ronnied

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

type Game struct {
	ID          string            `json:"id"`
	Name        string            `name:"name"`
	StartedAt   *time.Time        `json:"started_at"`
	Memberships []*GameMembership `json:"memberships"`
}

func (g *Game) String() string {
	return fmt.Sprintf("Game{ID: %s, Name: %s, StartedAt: %s}", g.ID, g.Name, g.StartedAt)
}

func (g *Game) MarshalGameString() string {
	output, err := json.Marshal(g)
	if err != nil {
		slog.Warn("error marshalling game", g.String(), err.Error())
		return ""
	}

	return string(output)
}

func UnmarshalGameString(input string) (*Game, error) {
	out := &Game{}
	err := json.Unmarshal([]byte(input), out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

type GameMembership struct {
	GameID   string `json:"game_id"`
	MemberID string `json:"member_id"`
}

type GameEntry struct {
	ID             string     `json:"id"`
	GameID         string     `json:"game_id"`
	MemberID       string     `json:"member_id"`
	Roll           int        `json:"roll"`
	AssignedTo     string     `json:"assigned_to"`
	Status         string     `json:"status"`
	CreatedDate    *time.Time `json:"created_date"`
	ReconciledDate *time.Time `json:"reconciled_date"`
}

type Tab struct {
	MemberID    string `json:"member_id"`
	GameEntryID string `json:"game_entry_id"`
}
