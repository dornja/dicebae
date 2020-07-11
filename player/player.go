package player

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"dicebae/baepi"
)

// PlayerHandler implements the BaeSayHandler interface for player character sheets.
type PlayerHandler struct {
	client           http.Client
	playerIDs        []int
	playerFirstNames []string
	nameToID         map[string]int
	charSheets       map[int]*CharacterSheet
}

func NewPlayerHandler(playerIDs []int) *PlayerHandler {
	ret := &PlayerHandler{
		client: http.Client{
			Timeout: 2 * time.Second,
		},
		playerIDs:  playerIDs,
		nameToID:   make(map[string]int),
		charSheets: make(map[int]*CharacterSheet),
	}

	for _, p := range ret.playerIDs {
		pj, err := ret.fetchPlayerJSON(p)
		if err != nil {
			// FIXME: this should just return an error and proceed or retry.
			panic("Failed to fetch player JSON:" + err.Error())
		}
		// Store the first name-part in lower case for easier matching. The
		// character sheet will have the properly capitalized form.
		// FIXME: This assumes characters have unique first names.
		n := strings.Split(strings.ToLower(pj.Name), " ")[0]
		ret.playerFirstNames = append(ret.playerFirstNames, n)
		ret.nameToID[n] = pj.ID
		ret.charSheets[pj.ID] = newCharacterSheet(pj)
	}
	sort.Strings(ret.playerFirstNames)
	for _, n := range ret.playerFirstNames {
		fmt.Printf("%s: %+v\n", n, ret.charSheets[ret.nameToID[n]])
	}

	return ret
}

func (ph *PlayerHandler) updateCharacterSheet(name string) error {
	id := ph.nameToID[name]
	js, err := ph.fetchPlayerJSON(id)
	if err != nil {
		return err
	}
	ph.charSheets[id] = newCharacterSheet(js)
	return nil
}

func (ph *PlayerHandler) ShouldSay(db baepi.DiceBae, e *baepi.Baevent) bool {
	if strings.HasPrefix(e.Message, "!who") {
		for _, n := range ph.playerFirstNames {
			if strings.Contains(strings.ToLower(e.Message), n) {
				return true
			}
		}
	}
	return false
}

func (ph *PlayerHandler) SayWithBae(db baepi.DiceBae, e *baepi.Baevent) (*baepi.Baesponse, error) {
	var resps []string
	for _, n := range ph.playerFirstNames {
		if strings.Contains(e.Message, n) {
			err := ph.updateCharacterSheet(n)
			if err != nil {
				resps = append(resps,
					fmt.Sprintf("**failed to update character sheet for %s, using cached version**", n))
			}
			cs := ph.charSheets[ph.nameToID[n]]
			resps = append(resps, cs.String())
		}
	}
	return &baepi.Baesponse{
		Message: strings.Join(resps, "\n"),
	}, nil
}
