package ronnied_actions

import (
	"context"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities/ronnied"
)

type CreateGameInput struct {
	Game *ronnied.Game
}

type CreateGameOutput struct {
	Game *ronnied.Game
}

type JoinGameInput struct {
	GameID   string
	MemberID string
}

type JoinGameOutput struct {
	Member *ronnied.GameMembership
}

type AddRollInput struct {
	GameID   string
	MemberID string
	Roll     int
}

type AddRollOutput struct {
	AssignedTo string
	Success    bool
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

type Interface interface {
	CreateGame(ctx context.Context, input *CreateGameInput) (*CreateGameOutput, error)
	JoinGame(ctx context.Context, input *JoinGameInput) (*JoinGameOutput, error)
	AddRoll(ctx context.Context, input *AddRollInput) (*AddRollOutput, error)
	GetTab(ctx context.Context, input *GetTabInput) (*GetTabOutput, error)
	PayDrink(ctx context.Context, input *PayDrinkInput) (*PayDrinkOutput, error)
}
