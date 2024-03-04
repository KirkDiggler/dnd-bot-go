package game

import (
	"github.com/KirkDiggler/dnd-bot-go/internal/entities/ronnied"
)

type CreateInput struct {
	Game *ronnied.Game
}

type CreateOutput struct {
	Game *ronnied.Game
}

type GetInput struct {
	ID string
}

type GetOutput struct {
	Game *ronnied.Game
}

type JoinInput struct {
	GameID   string
	MemberID string
}

type JoinOutput struct {
	Member *ronnied.GameMembership
}

type LeaveInput struct {
	GameID   string
	MemberID string
}

type LeaveOutput struct{}

type AddEntryInput struct {
	GameID     string
	MemberID   string
	Roll       int
	AssignedTo string
}

type AddEntryOutput struct {
	Entry *ronnied.GameEntry
}

type GetTabInput struct {
	GameID   string
	MemberID string
}

type GetTabOutput struct {
	Count int
}

type PayDrinkInput struct {
	GameID   string
	MemberID string
}

type PayDrinkOutput struct{}
