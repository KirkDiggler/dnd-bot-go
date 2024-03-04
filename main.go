package main

import (
	"flag"
	"github.com/KirkDiggler/dnd-bot-go/internal/managers/ronnied_actions"
	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/ronnied/game"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/choice"

	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/character_creation"

	"github.com/KirkDiggler/dnd-bot-go/internal/managers/characters"

	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/character"
	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/party"

	"github.com/KirkDiggler/dnd-bot-go/clients/dnd5e"

	"github.com/KirkDiggler/dnd-bot-go/discordbot"
)

var (
	token      string
	guildID    string
	appID      string
	redistHost string
)

func init() {
	flag.StringVar(&token, "token", "",
		"Bot token")
	flag.StringVar(&guildID, "guild", "",
		"Guild ID")
	flag.StringVar(&appID, "app", "",
		"Application ID")
	flag.StringVar(&redistHost, "redis", "localhost:6379",
		"Redis host")
	flag.Parse()

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
}

func main() {
	if token == "" || guildID == "" || appID == "" {
		flag.Usage()
		return
	}
	dnd5eClient, err := dnd5e.New(&dnd5e.Config{
		HttpClient: http.DefaultClient,
	})
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: redistHost,
	})

	partyRepo, err := party.New(&party.Config{
		Client: redisClient,
	})
	if err != nil {
		panic(err)
	}

	charRepo, err := character.New(&character.Config{
		Client: redisClient,
	})
	if err != nil {
		panic(err)
	}

	stateRepo, err := character_creation.New(&character_creation.Config{
		Client: redisClient,
	})
	if err != nil {
		panic(err)
	}

	choiceRepo, err := choice.New(&choice.Config{
		Client: redisClient,
	})
	if err != nil {
		panic(err)
	}

	charManager, err := characters.New(&characters.Config{
		Client:        dnd5eClient,
		CharacterRepo: charRepo,
		StateRepo:     stateRepo,
		ChoiceRepo:    choiceRepo,
	})
	if err != nil {
		panic(err)
	}

	gameRepo, err := game.NewRedis(&game.Config{
		Client: redisClient,
	})
	if err != nil {
		panic(err)
	}

	gameActions, err := ronnied_actions.NewManager(&ronnied_actions.ManagerConfig{
		GameRepo: gameRepo,
	})
	if err != nil {
		panic(err)
	}

	bot, err := discordbot.New(&discordbot.Config{
		Token:          token,
		GuildID:        guildID,
		AppID:          appID,
		DnD5EClient:    dnd5eClient,
		PartyRepo:      partyRepo,
		CharacterRepo:  charManager,
		RonnieDActions: gameActions,
	})
	if err != nil {
		panic(err)
	}

	err = bot.Start()
	if err != nil {
		panic(err)
	}

	defer func(bot discordbot.Interface) {
		err := bot.Close()
		if err != nil {
			panic(err)
		}
	}(bot)

	stchan := make(chan os.Signal, 1)
	signal.Notify(stchan, syscall.SIGTERM, os.Interrupt, syscall.SIGSEGV)

	for {
		select {
		case <-stchan:
			return
		default:
		}
		time.Sleep(time.Second)
	}

}
