package actions

import (
	"log"

	"github.com/ihrk/microbot/internal/bot"
	"github.com/ihrk/microbot/internal/config"
	"github.com/ihrk/microbot/internal/creds"
	"github.com/ihrk/microbot/internal/extra/riot"
)

var queueTypes = map[string]string{
	"solo": "RANKED_SOLO_5x5",
	"flex": "RANKED_FLEX_SR",
}

func Elo(cfg config.Settings) bot.Handler {
	region := cfg.MustString("region")
	summonerName := cfg.MustString("summonerName")
	queueType := cfg.StringWithDefault("queueType", "solo")

	tp, ok := queueTypes[queueType]
	if !ok {
		log.Fatalf("unknown lol queue type: %s\n", queueType)
	}

	c, err := riot.NewClient(region, creds.RiotAPIKey())
	if err != nil {
		log.Fatalf("riot client error: %v\n", err)
	}

	summoner, err := c.GetSummonerByName(summonerName)
	if err != nil {
		log.Fatalf("riot client error: %v\n", err)
	}

	return bot.HandlerFunc(func(s *bot.Sender) {
		entries, err := c.GetLeagueEntriesBySummoner(summoner.ID)
		if err != nil {
			log.Printf("riot api error: %v\n", err)
			return
		}

		var entry *riot.LeagueEntry

		for i := range entries {
			if entries[i].QueueType == tp {
				entry = &entries[i]
			}
		}

		if entry == nil {
			s.Reply("Rank not found")
			return
		}

		s.Reply(entry.String())
	})
}
