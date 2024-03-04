package monster

import (
	"context"
	"encoding/json"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"testing"
)

type monsterSuite struct {
	suite.Suite

	ctx         context.Context
	fixture     *Redis
	redisMock   redismock.ClientMock
	monster     *entities.Monster
	jsonMonster string
}

func (s *monsterSuite) SetupTest() {
	s.ctx = context.Background()
	client, redisMock := redismock.NewClientMock()
	s.fixture = &Redis{
		client: client,
	}
	s.redisMock = redisMock
	s.monster = &entities.Monster{
		Key: "1234",
		ID:  "5678",
	}
	buf, _ := json.Marshal(s.monster)
	s.jsonMonster = string(buf)
}

func (s *monsterSuite) TestGetMonster() {
	s.redisMock.ExpectGet("monster:" + s.monster.ID).SetVal(s.jsonMonster)

	result, err := s.fixture.GetMonster(s.ctx, s.monster.ID)
	s.NoError(err)
	s.NotNil(result)
	s.Equal(s.monster, result)
}

func (s *monsterSuite) TestGetMonsterNotFound() {
	s.redisMock.ExpectGet("monster:" + s.monster.ID).SetErr(redis.Nil)

	result, err := s.fixture.GetMonster(s.ctx, s.monster.ID)
	s.Error(err)
	s.Nil(result)
}

func (s *monsterSuite) TestGetMonsterError() {
	s.redisMock.ExpectGet("monster:" + s.monster.ID).SetErr(redis.Nil)

	result, err := s.fixture.GetMonster(s.ctx, s.monster.ID)
	s.Error(err)
	s.Nil(result)
}

func (s *monsterSuite) TestPutMonster() {
	s.redisMock.ExpectSet("monster:"+s.monster.ID, s.jsonMonster, 0).SetVal("OK")

	err := s.fixture.PutMonster(s.ctx, s.monster)
	s.NoError(err)
}

func (s *monsterSuite) TestPutMonsterError() {
	s.redisMock.ExpectSet("monster:"+s.monster.ID, s.jsonMonster, 0).SetErr(redis.Nil)

	err := s.fixture.PutMonster(s.ctx, s.monster)
	s.Error(err)
}

func TestMonster(t *testing.T) {
	suite.Run(t, new(monsterSuite))
}
