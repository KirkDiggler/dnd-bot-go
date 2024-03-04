package character_creation

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"

	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
)

type redisRepo struct {
	client redis.UniversalClient
}

type Config struct {
	Client redis.UniversalClient
}

func New(cfg *Config) (Repository, error) {
	if cfg == nil {
		return nil, dnderr.NewMissingParameterError("cfg")
	}

	if cfg.Client == nil {
		return nil, dnderr.NewMissingParameterError("cfg.Client")
	}

	return &redisRepo{
		client: cfg.Client,
	}, nil
}

func getCharacterCreationKey(id string) string {
	return fmt.Sprintf("character_creation:%s", id)
}

func (r *redisRepo) Get(ctx context.Context, id string) (*entities.CharacterCreation, error) {
	if id == "" {
		return nil, dnderr.NewMissingParameterError("id")
	}

	result := r.client.Get(ctx, getCharacterCreationKey(id))
	if result.Err() != nil {
		if result.Err() == redis.Nil {
			return nil, dnderr.NewNotFoundError("character creation not found")
		}

		return nil, result.Err()
	}

	return jsonToCharacterCreation(result.Val()), nil
}

func (r *redisRepo) Put(ctx context.Context, state *entities.CharacterCreation) (*entities.CharacterCreation, error) {
	if state == nil {
		return nil, dnderr.NewMissingParameterError("state")
	}

	result := r.client.Set(ctx, getCharacterCreationKey(state.CharacterID), characterCreateToJSON(state), 0)
	if result.Err() != nil {
		return nil, result.Err()
	}

	return state, nil
}
