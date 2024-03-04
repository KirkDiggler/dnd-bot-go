package party

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"

	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
	"github.com/KirkDiggler/dnd-bot-go/internal/types"

	"github.com/KirkDiggler/dnd-bot-go/dnderr"
)

type redisRepo struct {
	client redis.UniversalClient
	uuider types.UUIDGenerator
}

type Config struct {
	Client redis.UniversalClient
}

func New(cfg *Config) (Interface, error) {
	if cfg == nil {
		return nil, dnderr.NewMissingParameterError("cfg")
	}

	if cfg.Client == nil {
		return nil, dnderr.NewMissingParameterError("cfg.Client")
	}

	return &redisRepo{
		client: cfg.Client,
		uuider: &types.GoogleUUID{},
	}, nil
}

func (r *redisRepo) GetParty(ctx context.Context, token string) (*entities.Party, error) {
	if token == "" {
		return nil, dnderr.NewMissingParameterError("token")
	}

	result := r.client.Get(ctx, getPartyKey(token))
	if result.Err() != nil {
		if result.Err() == redis.Nil {
			return nil, dnderr.NewNotFoundError(fmt.Sprintf("token: %s not found", token))
		}

		return nil, result.Err()
	}

	return jsonToParty(result.Val()), nil
}
func getPartyKey(token string) string {
	return fmt.Sprintf("party:%s", token)
}

func (r *redisRepo) CreateParty(ctx context.Context, party *entities.Party) (*entities.Party, error) {
	if party == nil {
		return nil, dnderr.NewMissingParameterError("party")
	}

	if party.PartySize == 0 {
		return nil, dnderr.NewMissingParameterError("party.PartySize")
	}

	if party.Name == "" {
		return nil, dnderr.NewMissingParameterError("party.Name")
	}

	token := r.uuider.New()
	party.Token = token

	err := r.client.Set(ctx, getPartyKey(token), partyToJson(party), 0).Err()
	if err != nil {
		return nil, err
	}
	return party, nil
}
