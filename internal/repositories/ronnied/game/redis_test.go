package game

import (
	"context"
	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities/ronnied"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/suite"
	"testing"
)

type suiteRepo struct {
	suite.Suite

	ctx context.Context

	mockRedis redismock.ClientMock

	fixture *Redis
}

func (s *suiteRepo) SetupTest() {
	s.ctx = context.Background()
	client, mock := redismock.NewClientMock()
	s.mockRedis = mock

	fixture, err := NewRedis(&Config{
		Client: client,
	})
	s.Require().NoError(err)

	s.fixture = fixture
}

func (s *suiteRepo) TestGet_ValidatesInput() {
	testCases := []struct {
		name          string
		input         *GetInput
		expectedError error
	}{
		{
			name:          "Nil Input",
			input:         nil,
			expectedError: dnderr.NewMissingParameterError("input"),
		}, {
			name:          "Empty ID",
			input:         &GetInput{},
			expectedError: dnderr.NewMissingParameterError("input.ID"),
		}, {
			name: "Valid",
			input: &GetInput{
				ID: "123",
			},
			expectedError: nil,
		},
	}

	stubbedGame := &ronnied.Game{
		ID:   "123",
		Name: "Dispensing Liberty",
	}

	// we have one valid case so we will mock that call
	s.mockRedis.ExpectGet(getGameKey("123")).SetVal(stubbedGame.MarshalGameString())
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := s.fixture.Get(s.ctx, tc.input)
			if tc.expectedError == nil {
				s.Require().NoError(err)
				return
			}
			s.Require().Equal(tc.expectedError, err)
		})
	}
}

func TestRepo(t *testing.T) {
	suite.Run(t, new(suiteRepo))
}
