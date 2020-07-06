package roll

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	rollRegexp = regexp.MustCompile(`(\d*)\s*[dD](\d+)\s*([+-]\s*\d+)?`)

	maxAbsModifier   = 100
	maxDieSize       = 1000
	maxComputedRolls = 100000000
)

type rollTokens struct {
	multiplier string
	die        string
	modifier   string
}

// HasRollRequest returns whether a user's free-form text includes a dice roll.
func HasRollRequest(msg string) bool {
	return len(rollRegexp.FindAllStringSubmatch(msg, -1)) > 0
}

// ParseRequest returns whether a user's free-form text includes a dice roll.
func ParseRollRequest(msg string) (*RollRequest, error) {
	var specs []*RollSpec
	var errs []error
	var isTroll bool
	var trollMsg string
	for _, sub := range rollRegexp.FindAllStringSubmatch(msg, -1) {
		tk := &rollTokens{
			multiplier: sub[1],
			die:        sub[2],
			modifier:   sub[3],
		}
		r, err := parseRollTokens(tk)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		// Validate strings before parsing for weird input.
		switch {
		case r.Multiplier > maxComputedRolls:
			isTroll = true
			trollMsg = "I ain't got that many dice."
		case r.Die < 2:
			isTroll = true
			trollMsg = fmt.Sprintf("A %d-sided die is pointless, you ass.", r.Die)
		case r.Die > maxDieSize:
			isTroll = true
			trollMsg = fmt.Sprintf("A d%d is basically a sphere, wtf.", r.Die)
		case r.Modifier > maxAbsModifier:
			isTroll = true
			trollMsg = "You can't add that much to a modifier, that's unreasonable."
		case r.Modifier < -maxAbsModifier:
			isTroll = true
			trollMsg = "You can't subtract that much from a modifier, that's unreasonable."
		}
		specs = append(specs, r)
	}
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to parse rolls: %v", errs)
	}
	return &RollRequest{
		Timestamp: time.Now(),
		Specs:     specs,
		IsTroll:   isTroll,
		TrollMsg:  trollMsg,
	}, nil
}

func parseRollTokens(rt *rollTokens) (*RollSpec, error) {
	if rt == nil {
		return nil, fmt.Errorf("passed tokens are nil")
	}

	var errs []error
	mul, err := parseMultiplier(rt.multiplier)
	if err != nil {
		errs = append(errs, err)
	}
	die, err := parseDie(rt.die)
	if err != nil {
		errs = append(errs, err)
	}
	mod, err := parseModifier(rt.modifier)
	if err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to parse %+v: %v", *rt, errs)
	}

	return &RollSpec{
		Multiplier: mul,
		Die:        die,
		Modifier:   mod,
	}, nil
}

func parseMultiplier(mul string) (int, error) {
	if mul == "" {
		return 1, nil
	}
	ret, err := strconv.Atoi(mul)
	if err != nil {
		return 1, fmt.Errorf("failed to parse multiplier %q: %v", mul, err)
	}
	return ret, nil
}

func parseDie(die string) (int, error) {
	if die == "" {
		return 20, fmt.Errorf("missing a value for the die")
	}
	ret, err := strconv.ParseInt(die, 10, 64)
	if err != nil {
		return 20, fmt.Errorf("failed to parse die %q: %v", die, err)
	}
	return int(ret), nil
}

func parseModifier(modifier string) (int, error) {
	if strings.TrimSpace(modifier) == "" {
		return 0, nil
	}
	modifier = strings.ReplaceAll(modifier, " ", "")
	m, err := strconv.Atoi(modifier)
	if err != nil {
		return 0, fmt.Errorf("failed to parse modifier %q: %v", modifier, err)
	}
	return m, nil
}
