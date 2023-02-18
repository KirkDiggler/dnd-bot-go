package discordbot

import (
	"log"

	"github.com/KirkDiggler/dnd-bot-go/discordbot/components/character"

	"github.com/KirkDiggler/dnd-bot-go/internal/managers/characters"

	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/party"

	"github.com/KirkDiggler/dnd-bot-go/clients/dnd5e"
	"github.com/KirkDiggler/dnd-bot-go/discordbot/components"
	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/bwmarrin/discordgo"
)

type bot struct {
	session            *discordgo.Session
	guildID            string
	appID              string
	registeredCommands []*discordgo.ApplicationCommand

	partyRepo          party.Interface
	partyComponent     *components.Party
	characterComponent *character.Character
	ronnieDComponent   *components.RonnieD
}

type Config struct {
	Token         string
	GuildID       string
	AppID         string
	DnD5EClient   dnd5e.Client
	PartyRepo     party.Interface
	CharacterRepo characters.Manager
}

func New(cfg *Config) (*bot, error) {
	if cfg == nil {
		return nil, dnderr.NewMissingParameterError("cfg")
	}

	if cfg.DnD5EClient == nil {
		return nil, dnderr.NewMissingParameterError("cfg.Client")
	}

	if cfg.Token == "" {
		return nil, dnderr.NewMissingParameterError("cfg.Token")
	}

	if cfg.GuildID == "" {
		return nil, dnderr.NewMissingParameterError("cfg.GuildID")
	}

	if cfg.AppID == "" {
		return nil, dnderr.NewMissingParameterError("cfg.AppID")
	}

	if cfg.PartyRepo == nil {
		return nil, dnderr.NewMissingParameterError("cfg.PartyRepo")
	}

	session, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return nil, err
	}

	session.Identify.Intents |= discordgo.IntentGuildMembers
	session.Identify.Intents |= discordgo.IntentsGuilds
	session.Identify.Intents |= discordgo.IntentsGuildMessageReactions

	partyComponent, err := components.NewParty(&components.PartyConfig{
		Session:   session,
		PartyRepo: cfg.PartyRepo,
	})
	if err != nil {
		return nil, err
	}

	characterComponent, err := character.NewCharacter(&character.CharacterConfig{
		Client:           cfg.DnD5EClient,
		CharacterManager: cfg.CharacterRepo,
	})
	if err != nil {
		return nil, err
	}

	ronniedComponent, err := components.NewRonnieD()
	if err != nil {
		return nil, err
	}

	return &bot{
		session:            session,
		appID:              cfg.AppID,
		guildID:            cfg.GuildID,
		registeredCommands: make([]*discordgo.ApplicationCommand, 0),
		partyRepo:          cfg.PartyRepo,
		partyComponent:     partyComponent,
		characterComponent: characterComponent,
		ronnieDComponent:   ronniedComponent,
	}, nil
}

func (b *bot) Start() error {
	b.session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	// RonnieD commands
	b.session.AddHandler(b.ronnieDComponent.HandleMessageCreate)
	b.session.AddHandler(b.ronnieDComponent.HandleInteractionCreate)
	ronnieCmd := b.ronnieDComponent.GetApplicationCommand()

	_, err := b.session.ApplicationCommandCreate(b.appID, b.guildID, ronnieCmd)
	if err != nil {
		return err
	}

	b.registeredCommands = append(b.registeredCommands, ronnieCmd)

	// Party commands
	b.session.AddHandler(b.partyComponent.HandleInteractionCreate)

	partyCmd := b.partyComponent.GetApplicationCommand()
	_, err = b.session.ApplicationCommandCreate(b.appID, b.guildID, partyCmd)
	if err != nil {
		return err
	}

	b.registeredCommands = append(b.registeredCommands, partyCmd)

	// Character commands
	b.session.AddHandler(b.characterComponent.HandleInteractionCreate)
	charCmd := b.characterComponent.GetApplicationCommand()
	_, err = b.session.ApplicationCommandCreate(b.appID, b.guildID, charCmd)
	if err != nil {
		return err
	}

	b.registeredCommands = append(b.registeredCommands, charCmd)
	err = b.session.Open()
	if err != nil {
		return err
	}

	return nil
}

func (b *bot) Close() error {
	for _, v := range b.registeredCommands {
		log.Printf("Removing command '%v'...", v.Name)

		err := b.session.ApplicationCommandDelete(b.appID, b.guildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
	return b.session.Close()
}
