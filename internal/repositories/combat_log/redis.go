package combat_log

import (
	"context"
	"encoding/json"
	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal/types"
	"github.com/go-redis/redis/v9"
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
		client: cfg.Client,
		uuider: &types.GoogleUUID{},
	}, nil
}

func encounterCombatLogKey(encounterID string) string {
	return "encounterCombatLog:" + encounterID
}

func getCombatLogKey(id string) string {
	return "combatLog:" + id
}

func combatLogToJson(combatLog *Data) (string, error) {
	if combatLog == nil {
		return "", dnderr.NewMissingParameterError("combatLog")
	}

	buf, err := json.Marshal(combatLog)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func jsonToCombatLog(jsonStr string) (*Data, error) {
	if jsonStr == "" {
		return nil, dnderr.NewMissingParameterError("jsonStr")
	}

	var combatLog Data
	err := json.Unmarshal([]byte(jsonStr), &combatLog)
	if err != nil {
		return nil, err
	}

	return &combatLog, nil
}

func (r *Redis) Create(ctx context.Context, combatLog *Data) (*Data, error) {
	if combatLog == nil {
		return nil, dnderr.NewMissingParameterError("combatLog")
	}

	if combatLog.ID != "" {
		return nil, dnderr.NewInvalidEntityError("combatLog.ID")
	}

	if combatLog.EncounterID == "" {
		return nil, dnderr.NewMissingParameterError("combatLog.EncounterID")
	}

	combatLog.ID = r.uuider.New()

	now := r.timeClock.Now()

	combatLog.CreatedAt = now

	jsonStr, err := combatLogToJson(combatLog)
	if err != nil {
		return nil, err
	}

	combatLogCount, err := r.client.ZCard(ctx, encounterCombatLogKey(combatLog.EncounterID)).Result()
	if err != nil {
		return nil, err
	}

	pipe := r.client.TxPipeline()
	pipe.Set(ctx, getCombatLogKey(combatLog.ID), jsonStr, 0)
	pipe.ZAdd(ctx, encounterCombatLogKey(combatLog.EncounterID), redis.Z{
		Score:  float64(combatLogCount),
		Member: getCombatLogKey(combatLog.ID),
	})

	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	return combatLog, nil
}

func (r *Redis) Get(ctx context.Context, id string) (*Data, error) {
	if id == "" {
		return nil, dnderr.NewMissingParameterError("id")
	}

	jsonStr, err := r.client.Get(ctx, getCombatLogKey(id)).Result()
	if err != nil {
		return nil, err
	}

	combatLog, err := jsonToCombatLog(jsonStr)
	if err != nil {
		return nil, err
	}

	return combatLog, nil
}

func (r *Redis) ListByEncounter(ctx context.Context, input *ListByEncounterInput) ([]*Data, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.EncounterID == "" {
		return nil, dnderr.NewMissingParameterError("input.EncounterID")
	}

	var combatLogKeys []string
	var err error

	var start, stop int64

	start = input.Offset
	stop = input.Offset + input.Limit - 1

	if input.Reverse {
		combatLogKeys, err = r.client.ZRevRange(ctx, encounterCombatLogKey(input.EncounterID), start, stop).Result()
	} else {
		combatLogKeys, err = r.client.ZRange(ctx, encounterCombatLogKey(input.EncounterID), start, stop).Result()
	}

	combatLogs := make([]*Data, len(combatLogKeys))
	for i, combatLogKey := range combatLogKeys {
		combatLogs[i], err = r.Get(ctx, combatLogKey)
		if err != nil {
			return nil, err
		}
	}

	return combatLogs, nil
}
