package components

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const ronnieRollBack = "ronnie-roll-back"

type RonnieD struct {
	messageID string
}

func NewRonnieD() (*RonnieD, error) {
	return &RonnieD{}, nil
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
	if m.Content == "thanks ronnie" || m.Content == "thank's ronnie" || m.Content == "thanks ronnie d" || m.Content == "thank's ronnie d" {
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
}

func (c *RonnieD) HandleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		switch i.ApplicationCommandData().Name {
		case "ronnied":
			switch i.ApplicationCommandData().Options[0].Name {
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
			},
		},
	}
}
