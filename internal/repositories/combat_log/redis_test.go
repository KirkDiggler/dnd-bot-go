package combat_log

import (
	"context"
	"encoding/json"
	"github.com/KirkDiggler/dnd-bot-go/internal/types"
	"github.com/go-redis/redis/v9"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type combatLogSuite struct {
	suite.Suite

	ctx           context.Context
	redisMock     redismock.ClientMock
	mockUuider    *types.MockUUID
	fixture       *Redis
	timeMock      *types.MockClock
	combatLog     *Data
	jsonCombatLog string
}

func (s *combatLogSuite) SetupTest() {
	s.ctx = context.Background()
	client, redisMock := redismock.NewClientMock()
	s.mockUuider = &types.MockUUID{}
	s.timeMock = &types.MockClock{}
	s.fixture = &Redis{
		client:    client,
		uuider:    s.mockUuider,
		timeClock: s.timeMock,
	}
	s.redisMock = redisMock
	s.combatLog = &Data{
		ID:          "1234",
		EncounterID: "5678",
	}
	buf, _ := json.Marshal(s.combatLog)
	s.jsonCombatLog = string(buf)
}

func TestCombatLog(t *testing.T) {
	suite.Run(t, new(combatLogSuite))
}

func (s *combatLogSuite) TestCreateCombatLog() {
	s.mockUuider.On("New").Return(s.combatLog.ID)
	now := time.Now()
	s.timeMock.On("Now").Return(now)

	expected := &Data{
		ID:          s.combatLog.ID,
		EncounterID: s.combatLog.EncounterID,
		CreatedAt:   now,
		PlayerID:    s.combatLog.PlayerID,
		MonsterID:   s.combatLog.MonsterID,
		RoomID:      s.combatLog.RoomID,
		AttackRoll:  s.combatLog.AttackRoll,
		Type:        s.combatLog.Type,
	}

	expectedJson, _ := json.Marshal(expected)

	s.redisMock.ExpectZCard(encounterCombatLogKey(s.combatLog.EncounterID)).SetVal(42)

	s.redisMock.ExpectTxPipeline()
	s.redisMock.ExpectSet(getCombatLogKey(s.combatLog.ID), string(expectedJson), 0).SetVal(s.jsonCombatLog)

	s.redisMock.ExpectZAdd(encounterCombatLogKey(s.combatLog.EncounterID), redis.Z{
		Score:  float64(42),
		Member: s.combatLog.ID,
	}).SetVal(42)

	s.redisMock.ExpectTxPipelineExec()

	s.combatLog.ID = ""

	result, err := s.fixture.Create(s.ctx, s.combatLog)
	s.NoError(err)
	s.NotNil(result)
	s.Equal(expected, result)
}
