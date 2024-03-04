package components

import (
	"context"
	"fmt"
	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal/managers/ronnied_actions"
	"log"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const ronnieRollBack = "ronnie-roll-back"

type RonnieD struct {
	messageID string
	manager   ronnied_actions.Interface
}

type RonnieDConfig struct {
	Manager ronnied_actions.Interface
}

func NewRonnieD(cfg *RonnieDConfig) (*RonnieD, error) {
	if cfg == nil {
		return nil, dnderr.NewMissingParameterError("cfg")
	}

	if cfg.Manager == nil {
		return nil, dnderr.NewMissingParameterError("cfg.Manager")
	}

	return &RonnieD{
		manager: cfg.Manager,
	}, nil
}

func (c *RonnieD) RollBack(s *discordgo.Session, i *discordgo.InteractionCreate) {
	oldInteraction := &discordgo.Interaction{AppID: i.AppID, Token: c.messageID}
	err := s.InteractionResponseDelete(oldInteraction)
	if err != nil {
		log.Print(err)
	}

	c.RonnieRoll(s, i)
}

func (c *RonnieD) RonnieRoll(s *discordgo.Session, i *discordgo.InteractionCreate) {
	roll := rand.Intn(6) + 1
	msgBuilder := strings.Builder{}
	var response *discordgo.InteractionResponse
	c.messageID = i.Token

	if roll == 6 {
		msgBuilder.WriteString(fmt.Sprintf("%s rolled a Crit! Pass a drink", i.Member.User.Username))
		response = &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msgBuilder.String(),
				Components: []discordgo.MessageComponent{
					&discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							&discordgo.Button{
								Label:    "Roll it back",
								Style:    discordgo.SuccessButton,
								CustomID: ronnieRollBack,
								Emoji: discordgo.ComponentEmoji{
									Name: "🍺",
								},
							},
						},
					},
				},
			},
		}
	} else if roll == 1 {
		msgBuilder.WriteString(fmt.Sprintf("%s rolled a 1, that's a drink!", i.Member.User.Username))
		response = &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msgBuilder.String(),
				Components: []discordgo.MessageComponent{
					&discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							&discordgo.Button{
								Label:    "Roll it back",
								Style:    discordgo.DangerButton,
								CustomID: ronnieRollBack,
								Emoji: discordgo.ComponentEmoji{
									Name: "🍺",
								},
							},
						},
					},
				},
			},
		}
	} else {
		msgBuilder.WriteString(fmt.Sprintf("%s rolled a %d, try again", i.Member.User.Username, roll))
		response = &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msgBuilder.String(),
			},
		}
	}
	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Print(err)
	}
}

func (c *RonnieD) RonnieRollBack(s *discordgo.Session, i *discordgo.InteractionCreate) {
	msgBuilder := strings.Builder{}
	msgBuilder.WriteString(fmt.Sprintf("%s rolled it back!", i.Member.User.Username))
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: msgBuilder.String(),
		},
	}
	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Print(err)
	}
}

func (c *RonnieD) HandleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	grabBag := []string{
		"You know it Bro!",
		"Right back at you",
		"This guy get's it",
		"💯",
	}

	result := grabBag[rand.Intn(len(grabBag))]
	if m.Content == "thanks ronnie" ||
		m.Content == "Thank's Ronnie" ||
		m.Content == "Thank's ronnie" ||
		m.Content == "thank's Ronnie" ||
		m.Content == "thank's ronnie" ||
		m.Content == "thanks ronnie d" ||
		m.Content == "thank's ronnie d" {
		_, err := s.ChannelMessageSend(m.ChannelID, result)
		if err != nil {
			log.Print(err)
		}
	} else if m.Content == "tanks ronnie" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Get a load of this guy")
		if err != nil {
			log.Print(err)
		}
	} else if m.Content == "there it is" {
		_, err := s.ChannelMessageSend(m.ChannelID, "It's right there")
		if err != nil {
			log.Print(err)
		}
	} else if m.Content == "comon ronnie" {
		_, err := s.ChannelMessageSend(m.ChannelID, "You got this")
		if err != nil {
			log.Print(err)
		}
	}

	if m.Content == "how about you ronnie" || m.Content == "how about you ronnie d" {
		_, err := s.ChannelMessageSend(m.ChannelID, "I'm good Bro, thanks for asking")
		if err != nil {
			log.Print(err)
		}
	}

	if m.Content == "tanks ronnie" || m.Content == "tanks Ronnie" {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Get a load of %s. I'm not a tank, I'm a healer", m.Author.Username))
		if err != nil {
			log.Print(err)
		}
	}
}

func (c *RonnieD) HandleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		switch i.ApplicationCommandData().Name {
		case "ronnied":
			switch i.ApplicationCommandData().Options[0].Name {
			case "creategame":
				c.CreateGame(s, i)
			case "roll":
				c.RonnieRoll(s, i)
			case "advise":
				grabBag := []string{
					fmt.Sprintf("%s asked Ronnie D for advice, Ronnie D says: that's a drink", i.Member.User.Username),
					fmt.Sprintf("%s asked Ronnie D for advice, Ronnie D says: pass a drink", i.Member.User.Username),
					fmt.Sprintf("%s asked Ronnie D for advice, Ronnie D says: social!", i.Member.User.Username),
				}

				result := grabBag[rand.Intn(len(grabBag))]

				log.Println(result)

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: result,
					},
				})
				if err != nil {
					log.Print(err)
				}
			}
		}
	case discordgo.InteractionMessageComponent:
		switch i.MessageComponentData().CustomID {
		case ronnieRollBack:
			c.RollBack(s, i)
		}
	}
}

func (c *RonnieD) GetTab(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	if data.Options[0].Name == "gettab" {
		result, err := c.manager.GetTab(context.Background(), &ronnied_actions.GetTabInput{
			MemberID: i.Member.User.ID,
		})
		if err != nil {
			log.Print(err)
			return
		}

		msg := fmt.Sprintf("Your tab is %d", result.Count)
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		if err != nil {
			log.Print(err)
		}
	}
}

// AddResult adds a result to a game
// TODO: move to the roll command.  All rolls should be sent and based on the success response we will send a message
func (c *RonnieD) AddResult(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	if data.Options[0].Name == "addresult" {
		gameID := data.Options[0].Options[0].StringValue()
		roll := data.Options[0].Options[1].IntValue()
		result, err := c.manager.AddRoll(context.Background(), &ronnied_actions.AddRollInput{
			GameID:   gameID,
			MemberID: i.Member.User.ID,
			Roll:     int(roll),
		})
		if err != nil {
			log.Print(err)
			return
		}

		if result.Success {
			msg := fmt.Sprintf("Drink assigned to %s", result.AssignedTo)
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msg,
				},
			})
			if err != nil {
				log.Print(err)
			}
		}

	}
}

func (c *RonnieD) JoinGame(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	if data.Options[0].Name == "joingame" {
		gameID := i.ChannelID
		_, err := c.manager.JoinGame(context.Background(), &ronnied_actions.JoinGameInput{
			GameID:   gameID,
			MemberID: i.Member.User.ID,
		})
		if err != nil {
			log.Print(err)
			return
		}

		msg := fmt.Sprintf("You joined the game")
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		if err != nil {
			log.Print(err)
		}
	}
}

func (c *RonnieD) CreateGame(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	if data.Options[0].Name == "creategame" {
		gameName := data.Options[0].Options[0].StringValue()

		msg := fmt.Sprintf("Game %s created, ID: %d", gameName, rand.Intn(1000))
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		if err != nil {
			log.Print(err)
		}
	}
}

func (c *RonnieD) GetApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ronnied",
		Description: "Leverage RonnieD's wisdom",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "roll",
				Description: "roll em!",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			}, {
				Name:        "advise",
				Description: "what should I do RonnieD?",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			}, {
				Name:        "creategame",
				Description: "Create a game and get the game ID",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "name",
						Description: "Name of the game",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
					},
				},
			}, {
				Name:        "joingame",
				Description: "Join a game",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "gameid",
						Description: "ID of the game",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
					},
				},
			}, {
				Name:        "addresult",
				Description: "Add a result to a game",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "gameid",
						Description: "ID of the game",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
					},
				},
			}, {
				Name:        "gettab",
				Description: "Get your tab",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
		},
	}
}
