package choice

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"

	"github.com/KirkDiggler/dnd-bot-go/internal/entities"

	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal/types"
)

type redisRepo struct {
	client redis.UniversalClient
	uuider types.UUIDGenerator
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
		uuider: &types.GoogleUUID{},
	}, nil
}

func getChoiceKey(id string, choiceType entities.ChoiceType) string {
	return fmt.Sprintf("choice:%s:%s", id, choiceType)
}

func (r *redisRepo) Get(ctx context.Context, input *GetInput) (*GetOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.CharacterID == "" {
		return nil, dnderr.NewMissingParameterError("input.PlayerID")
	}

	if input.Type == "" {
		return nil, dnderr.NewMissingParameterError("input.Type")
	}

	result := r.client.Get(context.Background(), getChoiceKey(input.CharacterID, input.Type))
	if result.Err() != nil {
		if result.Err() == redis.Nil {
			return nil, dnderr.NewNotFoundError(fmt.Sprintf("choice not found: %s:%s", input.CharacterID, input.Type))
		}

		return nil, result.Err()
	}
	jsonResult := result.Val()

	data := jsonToData(jsonResult)
	return &GetOutput{
		CharacterID: data.CharacterID,
		Type:        typeToChoiceType(data.Type),
		Choices:     datasToChoices(data.Choices),
	}, nil
}

func (r *redisRepo) Put(ctx context.Context, input *PutInput) error {
	if input == nil {
		return dnderr.NewMissingParameterError("input")
	}

	if input.CharacterID == "" {
		return dnderr.NewMissingParameterError("input.PlayerID")
	}

	if input.Type == "" {
		return dnderr.NewMissingParameterError("input.Type")
	}

	if input.Choices == nil {
		return dnderr.NewMissingParameterError("input.Choices")
	}

	choices := choicesToDatas(input.Choices)
	data := &Data{
		CharacterID: input.CharacterID,
		Type:        choiceTypeToType(input.Type),
		Choices:     choices,
	}

	jsonData := dataToJSON(data)

	result := r.client.Set(context.Background(), getChoiceKey(input.CharacterID, input.Type), jsonData, 0)
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}
