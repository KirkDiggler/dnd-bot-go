package party

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
	"github.com/KirkDiggler/dnd-bot-go/internal/types"

	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

type partySuite struct {
	suite.Suite

	ctx         context.Context
	fixture     *redisRepo
	redisMock   redismock.ClientMock
	mockUuider  *types.MockUUID
	token       string
	party       *entities.Party
	jsonPayload string
}

func (s *partySuite) SetupTest() {
	s.ctx = context.Background()
	client, redisMock := redismock.NewClientMock()
	s.redisMock = redisMock
	s.mockUuider = &types.MockUUID{}
	s.party = &entities.Party{
		Name:      "Test Party",
		PartySize: 5,
		Token:     "1234",
	}
	s.token = "1234"
	jsonString, _ := json.Marshal(s.party)
	s.jsonPayload = string(jsonString)
	s.fixture = &redisRepo{
		client: client,
		uuider: s.mockUuider,
	}
}

func (s *partySuite) TestGetParty() {
	s.redisMock.ExpectGet(getPartyKey(s.token)).SetVal(s.jsonPayload)

	result, err := s.fixture.GetParty(s.ctx, s.token)
	s.NoError(err)
	s.NotNil(result)
	s.Equal(s.party, result)

}

func (s *partySuite) TestGetPartyNotFound() {
	s.redisMock.ExpectGet(getPartyKey(s.token)).SetErr(redis.Nil)

	result, err := s.fixture.GetParty(s.ctx, s.token)
	s.Error(err)
	s.Nil(result)
	s.EqualError(err, dnderr.NewNotFoundError("token: 1234 not found").Error())
}

func (s *partySuite) TestGetPartyError() {
	s.redisMock.ExpectGet(getPartyKey(s.token)).SetErr(errors.New("test error"))

	result, err := s.fixture.GetParty(s.ctx, s.token)
	s.Error(err)
	s.Nil(result)
	s.EqualError(err, errors.New("test error").Error())
}

func (s *partySuite) TestCreateParty() {
	s.redisMock.ExpectSet(getPartyKey(s.token), s.jsonPayload, 0).SetVal("OK")
	s.mockUuider.On("New").Return(s.token)

	result, err := s.fixture.CreateParty(s.ctx, &entities.Party{
		PartySize: 5,
		Name:      "Test Party",
	})

	s.NoError(err)
	s.NotNil(result)
	s.Equal(s.party, result)
}

func (s *partySuite) TestCreatePartyError() {
	s.redisMock.ExpectSet(getPartyKey(s.token), s.jsonPayload, 0).SetErr(errors.New("test error"))
	s.mockUuider.On("New").Return(s.token)

	result, err := s.fixture.CreateParty(s.ctx, &entities.Party{
		PartySize: 5,
		Name:      "Test Party",
	})

	s.Error(err)
	s.Nil(result)

	s.EqualError(err, errors.New("test error").Error())
}
func TestParty(t *testing.T) {
	suite.Run(t, new(partySuite))
}
