package player

import (
	"fmt"
)

type CharacterSheet struct {
	PlayerName string
	Class      string
	Level      int
	CurrentHP  int
	TotalHP    int

	Str int
	Dex int
	Con int
	Int int
	Wis int
	Cha int

	//TODO: Proficiency []string
	//TODO: Spell Slots []string
}

func fmtStat(stat int) string {
	return fmt.Sprintf("%d(**%+d**)", stat, statMod(stat))
}

func statMod(stat int) int {
	base := stat - 10
	if base > 0 {
		return base / 2
	} else {
		return (base - 1) / 2
	}
}

func (cs *CharacterSheet) String() string {

	stats := fmt.Sprintf("Str:%s Dex:%s Con:%s Int:%s Wis:%s Cha:%s",
		fmtStat(cs.Str), fmtStat(cs.Dex), fmtStat(cs.Con),
		fmtStat(cs.Int), fmtStat(cs.Wis), fmtStat(cs.Cha),
	)
	return fmt.Sprintf(
		"**%s:** Level %d %s, %d/%d HP\n%s", cs.PlayerName, cs.Level, cs.Class, cs.CurrentHP, cs.TotalHP, stats,
	)
}

func newCharacterSheet(p *DNDBeyondJSON) *CharacterSheet {
	statmap := make(map[int]int)
	for _, s := range p.Stats {
		statmap[s.ID] = s.Value
	}

	var pc PlayerClass
	if len(p.Classes) > 0 {
		pc = p.Classes[0]
	}

	cs := &CharacterSheet{
		PlayerName: p.Name,
		Class:      pc.Definition.Name,
		Level:      pc.Level,
		Str:        statmap[1],
		Dex:        statmap[2],
		Con:        statmap[3],
		Int:        statmap[4],
		Wis:        statmap[5],
		Cha:        statmap[6],
	}

	for _, mod := range append(p.Modifiers.Race, p.Modifiers.Class...) {
		switch mod.SubType {
		case "strength-score":
			cs.Str += mod.Value
		case "dexterity-score":
			cs.Dex += mod.Value
		case "constitution-score":
			cs.Con += mod.Value
		case "intelligence-score":
			cs.Int += mod.Value
		case "wisdom-score":
			cs.Wis += mod.Value
		case "charisma-score":
			cs.Cha += mod.Value
		}
	}

	hp := p.BaseHitPoints
	conMod := (cs.Con - 10) / 2
	hp += conMod * cs.Level

	cs.TotalHP = hp
	cs.CurrentHP = hp - p.RemovedHitPoints

	return cs
}
