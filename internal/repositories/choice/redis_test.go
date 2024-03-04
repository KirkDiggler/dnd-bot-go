package choice

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/KirkDiggler/dnd-bot-go/internal/entities"

	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/redis/go-redis/v9"

	"github.com/stretchr/testify/suite"

	"github.com/go-redis/redismock/v9"
)

type choiceSuite struct {
	suite.Suite

	ctx         context.Context
	fixture     *redisRepo
	redisMock   redismock.ClientMock
	characterID string
	choiceType  entities.ChoiceType
	data        *Data
	jsonPayload string
}

func (s *choiceSuite) SetupTest() {
	s.ctx = context.Background()
	client, redisMock := redismock.NewClientMock()
	s.redisMock = redisMock
	s.characterID = "1234"
	s.choiceType = entities.ChoiceTypeEquipment
	s.data = &Data{
		CharacterID: s.characterID,
		Type:        choiceTypeToType(s.choiceType),
		Choices: []*Choice{{
			Name:   "(a) a martial weapon and a shield or (b) two martial weapons",
			Count:  1,
			Type:   choiceTypeToType(s.choiceType),
			Status: StatusActive,
			Options: []*Option{{
				Multiple: &MultipleOption{
					Status: StatusActive,
					Options: []*Option{{
						Choice: &Choice{
							Name:   "a martial weapon",
							Count:  1,
							Type:   choiceTypeToType(s.choiceType),
							Status: StatusSelected,
							Options: []*Option{{
								Reference: &ReferenceOption{
									Reference: &ReferenceItem{
										Key:  "battleaxe",
										Name: "Battleaxe",
									},
								},
							}, {
								Reference: &ReferenceOption{
									Status: StatusSelected,
									Reference: &ReferenceItem{
										Key:  "flail",
										Name: "Flail",
									},
								},
							}},
						},
					}, {
						CountedReference: &CountedReferenceOption{
							Status: "selected",
							Count:  1,
							Reference: &ReferenceItem{
								Key:  "shield",
								Name: "Shield",
							},
						},
					}},
				},
			}, {
				Choice: &Choice{
					Name:  "two martial weapons",
					Count: 2,
					Type:  choiceTypeToType(s.choiceType),
					Options: []*Option{{
						Reference: &ReferenceOption{
							Reference: &ReferenceItem{
								Key:  "handaxe",
								Name: "Handaxe",
							},
						},
					}, {
						Reference: &ReferenceOption{
							Reference: &ReferenceItem{
								Key:  "longsword",
								Name: "Longsword",
							},
						},
					}},
				},
			}},
		}},
	}
	jsonString := dataToJSON(s.data)
	s.jsonPayload = jsonString
	s.fixture = &redisRepo{
		client: client,
	}
}

func (s *choiceSuite) TestGet() {
	s.redisMock.ExpectGet(getChoiceKey(s.characterID, s.choiceType)).SetVal(s.jsonPayload)
	actual, err := s.fixture.Get(s.ctx, &GetInput{
		CharacterID: s.characterID,
		Type:        s.choiceType,
	})

	s.NoError(err)
	s.Equal(&GetOutput{
		CharacterID: s.characterID,
		Type:        s.choiceType,
		Choices:     datasToChoices(s.data.Choices),
	}, actual)
}

func (s *choiceSuite) TestGetNotFound() {
	s.redisMock.ExpectGet(getChoiceKey(s.characterID, s.choiceType)).SetErr(redis.Nil)
	actual, err := s.fixture.Get(s.ctx, &GetInput{
		CharacterID: s.characterID,
		Type:        s.choiceType,
	})

	s.EqualError(err, dnderr.NewNotFoundError(fmt.Sprintf("choice not found: %s:%s", s.characterID, s.choiceType)).Error())
	s.Nil(actual)
}

func (s *choiceSuite) TestGetError() {
	s.redisMock.ExpectGet(getChoiceKey(s.characterID, s.choiceType)).SetErr(errors.New("some error"))
	actual, err := s.fixture.Get(s.ctx, &GetInput{
		CharacterID: s.characterID,
		Type:        s.choiceType,
	})

	s.EqualError(err, errors.New("some error").Error())
	s.Nil(actual)
}

func (s *choiceSuite) TestPut() {
	s.redisMock.ExpectSet(getChoiceKey(s.characterID, s.choiceType), s.jsonPayload, 0).SetVal("OK")
	err := s.fixture.Put(s.ctx, &PutInput{
		CharacterID: s.characterID,
		Type:        s.choiceType,
		Choices:     datasToChoices(s.data.Choices),
	})

	s.NoError(err)
}

func (s *choiceSuite) TestPutError() {
	s.redisMock.ExpectSet(getChoiceKey(s.characterID, s.choiceType), s.jsonPayload, 0).SetErr(errors.New("some error"))
	err := s.fixture.Put(s.ctx, &PutInput{
		CharacterID: s.characterID,
		Type:        s.choiceType,
		Choices:     datasToChoices(s.data.Choices),
	})

	s.EqualError(err, errors.New("some error").Error())
}

func TestChoice(t *testing.T) {
	suite.Run(t, new(choiceSuite))
}
