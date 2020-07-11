package player

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type DNDBeyondJSON struct {
	ID               int           `json:id`
	Name             string        `json:name`
	BaseHitPoints    int           `json:baseHitPoints`
	RemovedHitPoints int           `json:removedHitPoints`
	Stats            []PlayerStats `json:stats`
	Modifiers        Modifiers     `json:modifiers`
	Classes          []PlayerClass
}

type PlayerClass struct {
	Level      int                   `json:level`
	Definition PlayerClassDefinition `json:definition`
}

type PlayerClassDefinition struct {
	Name    string `json:name`
	HitDice int    `json:hitDice`
}

type PlayerStats struct {
	ID    int `json:id`
	Value int `json:value`
}

type Modifiers struct {
	Race  []Modifier `json:race`
	Class []Modifier `json:class`
}

type Modifier struct {
	ID          string `json:id`
	EntityID    int    `json:entityId`
	SubType     string `json:subType`
	TypeName    string `json:friendlyTypeName`
	SubTypeName string `json:friendlySubtypeName`
	Value       int    `json:value`
}

type SpellSlots struct {
	Level     int `json:level`
	Used      int `json:used`
	Available int `json:available`
}

func (ph *PlayerHandler) fetchPlayerJSON(playerID int) (*DNDBeyondJSON, error) {
	url := fmt.Sprintf("https://www.dndbeyond.com/character/%d/json", playerID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create an http request: %v", err)
	}
	req.Header.Set("User-Agent", "dicebae")

	res, err := ph.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %v", err)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading http result body failed: %v", err)
	}

	p := &DNDBeyondJSON{}
	if err := json.Unmarshal(body, p); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %v", err)
	}
	return p, nil
}
