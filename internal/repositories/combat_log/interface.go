package combat_log

import "context"

type Repository interface {
	Create(ctx context.Context, combatLog *Data) (*Data, error)
	Get(ctx context.Context, id string) (*Data, error)
	ListByEncounter(ctx context.Context, input *ListByEncounterInput) ([]*Data, error)
}
