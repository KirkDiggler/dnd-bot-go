package character

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"testing"

	"github.com/KirkDiggler/dnd-bot-go/internal/entities"

	"github.com/KirkDiggler/dnd-bot-go/internal/types"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/suite"
)

type characterSuite struct {
	suite.Suite

	ctx         context.Context
	fixture     *redisRepo
	redisMock   redismock.ClientMock
	mockUuider  *types.MockUUID
	id          string
	data        *Data
	character   *entities.Character
	jsonPayload string
}

func (s *characterSuite) SetupTest() {
	s.ctx = context.Background()
	client, redisMock := redismock.NewClientMock()
	s.redisMock = redisMock
	s.mockUuider = &types.MockUUID{}
	s.id = "1234"
	s.character = &entities.Character{
		ID:      s.id,
		OwnerID: s.id,
		Name:    "Test Character",
		Race: &entities.Race{
			Key: "elf",
		},
		Class: &entities.Class{
			Key: "fighter",
		},
	}
	s.data = characterToData(s.character)

	jsonString := dataToJSON(s.data)
	s.jsonPayload = jsonString
	s.fixture = &redisRepo{
		client: client,
		uuider: s.mockUuider,
	}
}

func (s *characterSuite) TestGetCharacter() {
	s.redisMock.ExpectGet(getCharacterKey(s.id)).SetVal(s.jsonPayload)

	result, err := s.fixture.Get(s.ctx, s.id)
	s.NoError(err)
	s.NotNil(result)
	s.Equal(s.data, result)
}

func (s *characterSuite) TestGetCharacterError() {
	s.redisMock.ExpectGet(getCharacterKey(s.id)).SetErr(errors.New("test error"))

	result, err := s.fixture.Get(s.ctx, s.id)
	s.Error(err)
	s.EqualError(err, "test error")
	s.Nil(result)
}
func (s *characterSuite) TestGetCharacterNotFound() {
	s.redisMock.ExpectGet(getCharacterKey(s.id)).SetErr(redis.Nil)

	result, err := s.fixture.Get(s.ctx, s.id)
	s.Error(err)
	s.EqualError(err, fmt.Sprintf("character id not found: %s", s.id))
	s.Nil(result)
}

func (s *characterSuite) TestCreateCharacter() {
	s.redisMock.ExpectSet(getCharacterKey(s.id), s.jsonPayload, 0).SetVal("OK")

	result, err := s.fixture.Put(s.ctx, s.character)
	s.NoError(err)
	s.NotNil(result)
	s.Equal(s.character, result)
}

func (s *characterSuite) TestCreateCharacterError() {
	s.redisMock.ExpectSet(getCharacterKey(s.id), s.jsonPayload, 0).SetErr(errors.New("test error"))

	result, err := s.fixture.Put(s.ctx, s.character)
	s.Error(err)
	s.EqualError(err, "test error")
	s.Nil(result)
}

func TestCharacter(t *testing.T) {
	suite.Run(t, new(characterSuite))
}
