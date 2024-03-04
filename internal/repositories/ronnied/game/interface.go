package game

import "context"

type Interface interface {
	Create(ctx context.Context, input *CreateInput) (*CreateOutput, error)
	Get(ctx context.Context, input *GetInput) (*GetOutput, error)
	Join(ctx context.Context, input *JoinInput) (*JoinOutput, error)
	AddEntry(ctx context.Context, input *AddEntryInput) (*AddEntryOutput, error)
	GetTab(ctx context.Context, input *GetTabInput) (*GetTabOutput, error)
	PayDrink(ctx context.Context, input *PayDrinkInput) (*PayDrinkOutput, error)
}
