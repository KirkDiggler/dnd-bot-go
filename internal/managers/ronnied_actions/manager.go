package ronnied_actions

import (
	"context"
	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/ronnied/game"
	"math/rand"
)

type Manager struct {
	gameRepo game.Interface
}

type ManagerConfig struct {
	GameRepo game.Interface
}

func NewManager(cfg *ManagerConfig) (*Manager, error) {
	if cfg == nil {
		return nil, dnderr.NewMissingParameterError("cfg")
	}

	if cfg.GameRepo == nil {
		return nil, dnderr.NewMissingParameterError("cfg.GameRepo")
	}

	return &Manager{
		gameRepo: cfg.GameRepo,
	}, nil
}

func (m *Manager) CreateGame(ctx context.Context, input *CreateGameInput) (*CreateGameOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.Game == nil {
		return nil, dnderr.NewMissingParameterError("input.Game")
	}

	if input.Game.Name == "" {
		return nil, dnderr.NewMissingParameterError("input.Game.Name")
	}

	result, err := m.gameRepo.Create(ctx, &game.CreateInput{
		Game: input.Game,
	})
	if err != nil {
		return nil, err
	}

	return &CreateGameOutput{
		Game: result.Game,
	}, nil
}

func (m *Manager) JoinGame(ctx context.Context, input *JoinGameInput) (*JoinGameOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.GameID == "" {
		return nil, dnderr.NewMissingParameterError("input.GameID")
	}

	if input.MemberID == "" {
		return nil, dnderr.NewMissingParameterError("input.MemberID")
	}

	result, err := m.gameRepo.Join(ctx, &game.JoinInput{
		GameID:   input.GameID,
		MemberID: input.MemberID,
	})
	if err != nil {
		return nil, err
	}

	return &JoinGameOutput{
		Member: result.Member,
	}, nil
}

func (m *Manager) AddRoll(ctx context.Context, input *AddRollInput) (*AddRollOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.GameID == "" {
		return nil, dnderr.NewMissingParameterError("input.GameID")
	}

	if input.MemberID == "" {
		return nil, dnderr.NewMissingParameterError("input.MemberID")
	}

	if shouldAddEntry(input.Roll) {
		var assignedTo = input.MemberID

		if input.Roll == 6 {
			gameResult, err := m.gameRepo.Get(ctx, &game.GetInput{
				ID: input.GameID,
			})
			if err != nil {
				return nil, err
			}

			// select a random membership for the game
			memberships := gameResult.Game.Memberships
			randIndex := rand.Intn(len(memberships))
			assignedTo = memberships[randIndex].MemberID

		}

		_, err := m.gameRepo.AddEntry(ctx, &game.AddEntryInput{
			GameID:     input.GameID,
			MemberID:   input.MemberID,
			Roll:       input.Roll,
			AssignedTo: assignedTo,
		})
		if err != nil {
			return nil, err
		}

		return &AddRollOutput{
			AssignedTo: assignedTo,
			Success:    true,
		}, nil
	}

	return &AddRollOutput{}, nil
}

func (m *Manager) GetTab(ctx context.Context, input *GetTabInput) (*GetTabOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.GameID == "" {
		return nil, dnderr.NewMissingParameterError("input.GameID")
	}

	if input.MemberID == "" {
		return nil, dnderr.NewMissingParameterError("input.MemberID")
	}

	result, err := m.gameRepo.GetTab(ctx, &game.GetTabInput{
		GameID:   input.GameID,
		MemberID: input.MemberID,
	})
	if err != nil {
		return nil, err
	}

	return &GetTabOutput{
		Count: result.Count,
	}, nil
}

func (m *Manager) PayDrink(ctx context.Context, input *PayDrinkInput) (*PayDrinkOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.GameID == "" {
		return nil, dnderr.NewMissingParameterError("input.GameID")
	}

	if input.MemberID == "" {
		return nil, dnderr.NewMissingParameterError("input.MemberID")
	}

	_, err := m.gameRepo.PayDrink(ctx, &game.PayDrinkInput{
		GameID:   input.GameID,
		MemberID: input.MemberID,
	})
	if err != nil {
		return nil, err
	}

	return &PayDrinkOutput{}, nil

}
func shouldAddEntry(roll int) bool {
	return roll == 1 || roll == 6
}
