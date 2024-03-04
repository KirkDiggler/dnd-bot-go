package room

import (
	"context"
	"encoding/json"
	"github.com/KirkDiggler/dnd-bot-go/dnderr"
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

func characterRoomKey(characterID string) string {
	return "characterRoom:" + characterID
}

func getRoomKey(id string) string {
	return "room:" + id
}

func roomToJson(room *Data) (string, error) {
	if room == nil {
		return "", dnderr.NewMissingParameterError("room")
	}

	buf, err := json.Marshal(room)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func jsonToRoom(jsonStr string) (*Data, error) {
	if jsonStr == "" {
		return nil, dnderr.NewMissingParameterError("jsonStr")
	}

	var room Data
	err := json.Unmarshal([]byte(jsonStr), &room)
	if err != nil {
		return nil, err
	}

	return &room, nil
}

// Create creates a room and assigns it to the owners index
func (r *Redis) Create(ctx context.Context, room *Data) (*Data, error) {
	if room == nil {
		return nil, dnderr.NewMissingParameterError("room")
	}

	if room.ID != "" {
		return nil, dnderr.NewInvalidEntityError("room.ID must be empty")
	}

	room.ID = r.uuider.New()

	jsonStr, err := roomToJson(room)
	if err != nil {
		return nil, err
	}

	roomCount, err := r.client.ZCard(ctx, characterRoomKey(room.PlayerID)).Result()
	if err != nil {
		return nil, err
	}

	// Use pipe to set the room data and add member to character rooms index
	pipe := r.client.TxPipeline()
	pipe.Set(ctx, getRoomKey(room.ID), jsonStr, 0)
	pipe.ZAdd(ctx, characterRoomKey(room.PlayerID), redis.Z{
		Score:  float64(roomCount),
		Member: getRoomKey(room.ID),
	})

	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (r *Redis) Update(ctx context.Context, room *Data) (*Data, error) {
	if room == nil {
		return nil, dnderr.NewMissingParameterError("room")
	}

	if room.ID == "" {
		return nil, dnderr.NewInvalidEntityError("room.ID must not be empty")
	}

	jsonStr, err := roomToJson(room)
	if err != nil {
		return nil, err
	}

	// TODO: verify room exists

	err = r.client.Set(ctx, getRoomKey(room.ID), jsonStr, 0).Err()
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (r *Redis) Get(ctx context.Context, id string) (*Data, error) {
	if id == "" {
		return nil, dnderr.NewMissingParameterError("id")
	}

	result := r.client.Get(ctx, getRoomKey(id))
	if result.Err() != nil {
		if result.Err() == redis.Nil {
			return nil, dnderr.NewNotFoundError("room not found")
		}

		return nil, result.Err()
	}

	return jsonToRoom(result.Val())
}

func (r *Redis) ListByPlayer(ctx context.Context, input *ListByPlayerInput) ([]*Data, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.PlayerID == "" {
		return nil, dnderr.NewMissingParameterError("input.PlayerID")
	}

	if input.Limit == 0 {
		input.Limit = 10
	}

	var roomKeys []string
	var err error

	var start, stop int64

	start = input.Offset
	stop = input.Offset + input.Limit - 1

	if input.Reverse {
		roomKeys, err = r.client.ZRevRange(ctx, characterRoomKey(input.PlayerID), start, stop).Result()
	} else {
		roomKeys, err = r.client.ZRange(ctx, characterRoomKey(input.PlayerID), start, stop).Result()
	}

	if err != nil {
		return nil, err
	}

	rooms := make([]*Data, len(roomKeys))
	for i, roomKey := range roomKeys {
		room, err := r.Get(ctx, roomKey)
		if err != nil {
			return nil, err
		}

		rooms[i] = room
	}

	return rooms, nil
}
