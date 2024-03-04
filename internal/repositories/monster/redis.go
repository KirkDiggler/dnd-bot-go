package monster

import (
	"context"
	"encoding/json"
	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
	"github.com/KirkDiggler/dnd-bot-go/internal/types"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client redis.UniversalClient
	uuider types.UUIDGenerator
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
		client: cfg.Client,
		uuider: &types.GoogleUUID{},
	}, nil
}

func getMonsterKey(ID string) string {
	return "monster:" + ID
}

func (r *Redis) GetMonster(ctx context.Context, ID string) (*entities.Monster, error) {
	if ID == "" {
		return nil, dnderr.NewMissingParameterError("ID")
	}

	result := r.client.Get(ctx, getMonsterKey(ID))
	if result.Err() != nil {
		if result.Err() == redis.Nil {
			return nil, dnderr.NewNotFoundError("monster not found")
		}

		return nil, result.Err()
	}

	return jsonToMonster(result.Val())
}

func jsonToMonster(input string) (*entities.Monster, error) {
	out := &entities.Monster{}

	err := json.Unmarshal([]byte(input), out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func monsterToJson(input *entities.Monster) (string, error) {
	out, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func (r *Redis) PutMonster(ctx context.Context, monster *entities.Monster) error {
	if monster == nil {
		return dnderr.NewMissingParameterError("monster")
	}

	key := getMonsterKey(monster.ID)
	jsonValue, err := monsterToJson(monster)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, jsonValue, 0).Err()
}
