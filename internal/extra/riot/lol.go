package riot

import "fmt"

type Summoner struct {
	ID        string `json:"id"`
	AccountID string `json:"accountID"`
}

func (c *Client) GetSummonerByName(name string) (*Summoner, error) {
	var s Summoner

	path := fmt.Sprintf("/lol/summoner/v4/summoners/by-name/%s", name)

	err := c.doRequest(path, &s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

type LeagueEntry struct {
	QueueType    string
	Tier         string
	Rank         string
	SummonerName string
	LeaguePoints int
	Wins         int
	Losses       int
}

func (e *LeagueEntry) String() string {
	return fmt.Sprintf("%s: %s %s %d lp, w/l: %d/%d",
		e.SummonerName,
		e.Tier,
		e.Rank,
		e.LeaguePoints,
		e.Wins,
		e.Losses,
	)
}

func (c *Client) GetLeagueEntriesBySummoner(summonerID string) ([]LeagueEntry, error) {
	var entries []LeagueEntry

	path := fmt.Sprintf("/lol/league/v4/entries/by-summoner/%s", summonerID)

	err := c.doRequest(path, &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}
