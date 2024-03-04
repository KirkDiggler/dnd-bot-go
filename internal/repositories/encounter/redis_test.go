package encounter

import (
	"context"
	"encoding/json"
	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal/types"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type encounterSuite struct {
	suite.Suite

	ctx           context.Context
	fixture       *Redis
	redisMock     redismock.ClientMock
	encounter     *Data
	jsonEncounter string
	uuiderMock    *types.MockUUID
	timeMock      *types.MockClock
}

func (s *encounterSuite) SetupTest() {
	s.ctx = context.Background()
	client, redisMock := redismock.NewClientMock()
	s.uuiderMock = &types.MockUUID{}
	s.timeMock = &types.MockClock{}
	s.fixture = &Redis{
		client:    client,
		uuider:    s.uuiderMock,
		timeClock: s.timeMock,
	}
	s.redisMock = redisMock
	s.encounter = &Data{
		ID:       "1234",
		PlayerID: "5678",
	}
	buf, _ := json.Marshal(s.encounter)
	s.jsonEncounter = string(buf)
}

func TestEncounter(t *testing.T) {
	suite.Run(t, new(encounterSuite))
}

func (s *encounterSuite) TestCreateEncounter() {
	s.uuiderMock.On("New").Return(s.encounter.ID)
	now := time.Now()
	s.timeMock.On("Now").Return(now)

	expected := &Data{
		ID:        s.encounter.ID,
		PlayerID:  s.encounter.PlayerID,
		CreatedAt: now,
		UpdatedAt: now,
		StartDate: s.encounter.StartDate,
		EndDate:   s.encounter.EndDate,
		Status:    s.encounter.Status,
		MonsterID: s.encounter.MonsterID,
		RoomID:    s.encounter.RoomID,
	}

	expectedJson, _ := json.Marshal(expected)

	s.redisMock.ExpectZCard(characterEncounterKey(s.encounter.PlayerID)).SetVal(42)

	s.redisMock.ExpectTxPipeline()
	s.redisMock.ExpectSet(getEncounterKey(s.encounter.ID), string(expectedJson), 0).SetVal(s.jsonEncounter)

	s.redisMock.ExpectZAdd(characterEncounterKey(s.encounter.PlayerID), redis.Z{
		Score:  43,
		Member: getEncounterKey(s.encounter.ID),
	}).SetVal(42)

	s.redisMock.ExpectTxPipelineExec()

	s.encounter.ID = ""

	result, err := s.fixture.Create(s.ctx, s.encounter)
	s.NoError(err)
	s.NotNil(result)
	s.Equal(expected, result)
}

func (s *encounterSuite) TestCreateEncounter_InvalidInput() {
	_, err := s.fixture.Create(s.ctx, nil)
	s.Error(dnderr.NewMissingParameterError("encounter"))
	s.EqualError(err, dnderr.NewMissingParameterError("encounter").Error())

	_, err = s.fixture.Create(s.ctx, &Data{})
	s.Error(err)
	s.EqualError(err, dnderr.NewMissingParameterError("encounter.PlayerID").Error())

	_, err = s.fixture.Create(s.ctx, &Data{
		ID: "1234",
	})
	s.Error(err)
	s.EqualError(err, dnderr.NewInvalidEntityError("encounter.ID").Error())
}

func (s *encounterSuite) TestGetEncounter() {
	s.redisMock.ExpectGet(getEncounterKey(s.encounter.ID)).SetVal(s.jsonEncounter)

	result, err := s.fixture.Get(s.ctx, s.encounter.ID)
	s.NoError(err)
	s.NotNil(result)
	s.Equal(s.encounter, result)
}

func (s *encounterSuite) TestGetEncounterNotFound() {
	s.redisMock.ExpectGet(getEncounterKey(s.encounter.ID)).RedisNil()

	result, _ := s.fixture.Get(s.ctx, s.encounter.ID)
	s.NoError(nil)
	s.Nil(result)
}

func (s *encounterSuite) TestUpdate() {
	now := time.Now()
	s.timeMock.On("Now").Return(now)

	existing := &Data{
		ID:        s.encounter.ID,
		PlayerID:  s.encounter.PlayerID,
		CreatedAt: now.Add(-1 * time.Hour),
		UpdatedAt: now.Add(-1 * time.Hour),
		StartDate: s.encounter.StartDate,
		EndDate:   s.encounter.EndDate,
		Status:    s.encounter.Status,
		MonsterID: s.encounter.MonsterID,
		RoomID:    s.encounter.RoomID,
	}

	expected := &Data{
		ID:        s.encounter.ID,
		PlayerID:  s.encounter.PlayerID,
		CreatedAt: now.Add(-1 * time.Hour),
		UpdatedAt: now,
		StartDate: s.encounter.StartDate,
		EndDate:   s.encounter.EndDate,
		Status:    s.encounter.Status,
		MonsterID: s.encounter.MonsterID,
		RoomID:    s.encounter.RoomID,
	}

	existingJson, _ := json.Marshal(existing)
	expectedJson, _ := json.Marshal(expected)

	s.redisMock.ExpectGet(getEncounterKey(s.encounter.ID)).SetVal(string(existingJson))
	s.redisMock.ExpectSet(getEncounterKey(s.encounter.ID), string(expectedJson), 0).SetVal(s.jsonEncounter)

	err := s.fixture.Update(s.ctx, s.encounter)
	s.NoError(err)
}

func (s *encounterSuite) TestListByPlayer() {
	s.redisMock.ExpectZRevRange(characterEncounterKey(s.encounter.PlayerID), 0, 10).SetVal([]string{
		getEncounterKey(s.encounter.ID),
	})

	s.redisMock.ExpectGet(getEncounterKey(s.encounter.ID)).SetVal(s.jsonEncounter)

	result, err := s.fixture.ListByPlayer(s.ctx, s.encounter.PlayerID)
	s.NoError(err)
	s.NotNil(result)
	s.Equal([]*Data{s.encounter}, result)
}
