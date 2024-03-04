package game

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities/ronnied"
	"github.com/KirkDiggler/dnd-bot-go/internal/types"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"strconv"
	"time"
)

type Redis struct {
	client redis.UniversalClient
	uuider types.UUIDGenerator
}

type Config struct {
	Client redis.UniversalClient
}

func NewRedis(cfg *Config) (*Redis, error) {
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

func (r *Redis) Get(ctx context.Context, input *GetInput) (*GetOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.ID == "" {
		return nil, dnderr.NewMissingParameterError("input.ID")
	}

	gameKey := getGameKey(input.ID)

	gameJson, err := r.client.Get(ctx, gameKey).Result()
	if err != nil {
		return nil, err
	}

	game, err := ronnied.UnmarshalGameString(gameJson)
	if err != nil {
		return nil, err
	}

	slog.Info("Get ", "game", game.String())

	return &GetOutput{
		Game: game,
	}, nil
}

func (r *Redis) Create(ctx context.Context, input *CreateInput) (*CreateOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.Game == nil {
		return nil, dnderr.NewMissingParameterError("input.Game")
	}

	slog.Info("Creating ", "game", input.Game.String())

	gameJson, err := json.Marshal(input.Game)
	if err != nil {
		return nil, err
	}

	gameKey := getGameKey(input.Game.ID)

	err = r.client.Set(ctx, gameKey, gameJson, 0).Err()
	if err != nil {
		return nil, err
	}

	return &CreateOutput{
		Game: input.Game,
	}, nil
}
func (r *Redis) Leave(ctx context.Context, input *LeaveInput) (*LeaveOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.GameID == "" {
		return nil, dnderr.NewMissingParameterError("input.GameID")
	}

	if input.MemberID == "" {
		return nil, dnderr.NewMissingParameterError("input.MemberID")
	}

	membershipKey := getGameMembershipKey(input.GameID)

	err := r.client.SRem(ctx, membershipKey, input.MemberID).Err()
	if err != nil {
		return nil, err
	}

	return &LeaveOutput{}, nil
}

func (r *Redis) Join(ctx context.Context, input *JoinInput) (*JoinOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.GameID == "" {
		return nil, dnderr.NewMissingParameterError("input.GameID")
	}

	gameKey := getGameKey(input.GameID)

	// make sure we have the game saved already
	gameJson, err := r.client.Get(ctx, gameKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, dnderr.NewNotFoundError("game " + input.GameID + " not found")
		}

		return nil, err
	}

	game, err := ronnied.UnmarshalGameString(gameJson)
	if err != nil {
		return nil, err
	}

	// we have a game, let's add the user to the game
	membershipKey := getGameMembershipKey(input.GameID)

	err = r.client.SAdd(ctx, membershipKey, input.MemberID).Err()
	if err != nil {
		return nil, err
	}

	return &JoinOutput{
		Member: &ronnied.GameMembership{
			GameID:   game.ID,
			MemberID: input.MemberID,
		},
	}, nil
}

func (r *Redis) AddEntry(ctx context.Context, input *AddEntryInput) (*AddEntryOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.GameID == "" {
		return nil, dnderr.NewMissingParameterError("input.GameID")
	}

	if input.MemberID == "" {
		return nil, dnderr.NewMissingParameterError("input.MemberID")
	}

	if input.AssignedTo == "" {
		return nil, dnderr.NewMissingParameterError("input.AssignedTo")
	}

	if input.Roll != 1 && input.Roll != 6 {
		return nil, dnderr.NewInvalidParameterError(strconv.Itoa(input.Roll), "input.Roll must be 1 or 6")
	}

	if input.Roll == 1 {
		input.AssignedTo = input.MemberID
	}

	// TODO: make sure we have a game with that id saved
	// TODO: ensure the assigned to and member have membership to the game
	entry := &ronnied.GameEntry{
		ID:       r.uuider.New(),
		GameID:   input.GameID,
		MemberID: input.MemberID,
		Roll:     input.Roll,
	}

	entryJson, err := json.Marshal(entry)
	if err != nil {
		return nil, err
	}

	entryKey := getGameEntryKey(entry.GameID)

	tx := r.client.TxPipeline()
	tx.Set(ctx, entryKey, entryJson, 0)
	tx.ZAdd(ctx, getGameListEntriesKey(entry.GameID), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: entry.ID,
	})
	tx.RPush(ctx, getMemberTabKey(input.GameID, entry.MemberID), entry.ID)

	_, err = tx.Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &AddEntryOutput{
		Entry: entry,
	}, nil
}

func (r *Redis) GetTab(ctx context.Context, input *GetTabInput) (*GetTabOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.MemberID == "" {
		return nil, dnderr.NewMissingParameterError("input.MemberID")
	}

	tabKey := getMemberTabKey(input.GameID, input.MemberID)

	count, err := r.client.LLen(ctx, tabKey).Result()
	if err != nil {
		return nil, err
	}

	return &GetTabOutput{
		Count: int(count),
	}, nil
}

func (r *Redis) PayDrink(ctx context.Context, input *PayDrinkInput) (*PayDrinkOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.GameID == "" {
		return nil, dnderr.NewMissingParameterError("input.GameID")
	}

	if input.MemberID == "" {
		return nil, dnderr.NewMissingParameterError("input.MemberID")
	}

	tabKey := getMemberTabKey(input.GameID, input.MemberID)

	tx := r.client.TxPipeline()
	tx.LPop(ctx, tabKey)
	tx.RPush(ctx, getMemberPaidTabKey(input.GameID, input.MemberID), input.MemberID)

	_, err := tx.Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &PayDrinkOutput{}, nil
}
func getGameMembershipKey(gameID string) string {
	return "game:members:" + gameID
}
func getGameEntryKey(gameID string) string {
	return "game:entry:" + gameID
}
func getGameListEntriesKey(gameID string) string {
	return "game:entries:" + gameID
}
func getMemberTabKey(gameID, memberID string) string {
	return "tab:game" + gameID + ":member:" + memberID
}
func getMemberPaidTabKey(gameID, memberID string) string {
	return "paid:game:" + gameID + ":member" + memberID
}

func getGameKey(id string) string {
	return "game:" + id
}
