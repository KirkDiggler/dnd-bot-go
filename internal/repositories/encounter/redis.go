package encounter

import (
	"context"
	"encoding/json"
	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal/types"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client    redis.UniversalClient
	uuider    types.UUIDGenerator
	timeClock types.TimeClock
}

type RedisConfig struct {
	Client redis.UniversalClient
}

func NewRedis(cfg *RedisConfig) (*Redis, error) {
	if cfg == nil {
		return nil, dnderr.NewMissingParameterError("cfg")
	}

	if cfg.Client == nil {
		return nil, dnderr.NewMissingParameterError("cfg.Client")
	}

	return &Redis{
		client:    cfg.Client,
		uuider:    &types.GoogleUUID{},
		timeClock: &types.Clock{},
	}, nil
}

func characterEncounterKey(characterID string) string {
	return "characterEncounter:" + characterID
}

func getEncounterKey(id string) string {
	return "encounter:" + id
}

func encounterToJson(encounter *Data) (string, error) {
	if encounter == nil {
		return "", dnderr.NewMissingParameterError("encounter")
	}

	buf, err := json.Marshal(encounter)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func jsonToEncounter(jsonStr string) (*Data, error) {
	if jsonStr == "" {
		return nil, dnderr.NewMissingParameterError("jsonStr")
	}

	var encounter Data
	err := json.Unmarshal([]byte(jsonStr), &encounter)
	if err != nil {
		return nil, err
	}

	return &encounter, nil
}

func (r *Redis) Create(ctx context.Context, encounter *Data) (*Data, error) {
	if encounter == nil {
		return nil, dnderr.NewMissingParameterError("encounter")
	}

	if encounter.ID != "" {
		return nil, dnderr.NewInvalidEntityError("encounter.ID")
	}

	if encounter.PlayerID == "" {
		return nil, dnderr.NewMissingParameterError("encounter.PlayerID")
	}

	encounter.ID = r.uuider.New()

	now := r.timeClock.Now()

	encounter.CreatedAt = now
	encounter.UpdatedAt = now

	jsonStr, err := encounterToJson(encounter)
	if err != nil {
		return nil, err
	}

	encounterCount, err := r.client.ZCard(ctx, characterEncounterKey(encounter.PlayerID)).Result()
	if err != nil {
		return nil, err
	}

	// Add encounter and by-character index with transaction
	pipe := r.client.TxPipeline()
	pipe.Set(ctx, getEncounterKey(encounter.ID), jsonStr, 0)
	pipe.ZAdd(ctx, characterEncounterKey(encounter.PlayerID), redis.Z{
		Score:  float64(encounterCount + 1),
		Member: getEncounterKey(encounter.ID),
	})

	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	return encounter, nil
}

func (r *Redis) Update(ctx context.Context, encounter *Data) error {
	if encounter == nil {
		return dnderr.NewMissingParameterError("encounter")
	}

	if encounter.ID == "" {
		return dnderr.NewInvalidEntityError("encounter.ID cannot be null")
	}

	existing, err := r.Get(ctx, encounter.ID)
	if err != nil {
		return err
	}

	encounter.CreatedAt = existing.CreatedAt
	encounter.UpdatedAt = r.timeClock.Now()

	jsonStr, err := encounterToJson(encounter)
	if err != nil {
		return err
	}

	err = r.client.Set(ctx, getEncounterKey(encounter.ID), jsonStr, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) Get(ctx context.Context, id string) (*Data, error) {
	if id == "" {
		return nil, dnderr.NewMissingParameterError("id")
	}

	jsonStr, err := r.client.Get(ctx, getEncounterKey(id)).Result()
	if err != nil {
		return nil, err
	}

	encounter, err := jsonToEncounter(jsonStr)
	if err != nil {
		return nil, err
	}

	return encounter, nil
}

func (r *Redis) ListByPlayer(ctx context.Context, playerID string) ([]*Data, error) {
	if playerID == "" {
		return nil, dnderr.NewMissingParameterError("playerID")
	}

	encounterKeys, err := r.client.ZRevRange(ctx, characterEncounterKey(playerID), 0, 10).Result()
	if err != nil {
		return nil, err
	}

	var encounters []*Data
	for _, encounterKey := range encounterKeys {
		jsonStr, err := r.client.Get(ctx, encounterKey).Result()
		if err != nil {
			return nil, err
		}

		encounter, err := jsonToEncounter(jsonStr)
		if err != nil {
			return nil, err
		}

		encounters = append(encounters, encounter)
	}

	return encounters, nil
}
