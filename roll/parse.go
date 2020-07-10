package roll

import (
	"fmt"
	//	"math/rand"
	"strconv"
	"strings"
)

func parseRollRequests(msg string) ([]*RollRequest, error) {
	var ret []*RollRequest
	var errs []error
	for _, sub := range rollRegexp.FindAllStringSubmatch(msg, -1) {
		mulStr, dieStr, modStr := sub[1], sub[2], sub[3]

		var perrs []string
		mul, err := parseMultiplier(mulStr)
		if err != nil {
			perrs = append(perrs, err.Error())
		}
		die, err := parseDie(dieStr)
		if err != nil {
			perrs = append(perrs, err.Error())
		}
		mod, err := parseModifier(modStr)
		if err != nil {
			perrs = append(perrs, err.Error())
		}
		if len(perrs) > 0 {
			errs = append(errs, fmt.Errorf(
				"failed to parse, mul:%q, die:%q, mod:%q: %q",
				mulStr, dieStr, modStr, strings.Join(perrs, ", "),
			))
		}
		r := &RollRequest{
			Multiplier: mul,
			Die:        die,
			Modifier:   mod,
		}
		r.TrollMsg = checkForTrolls(r)

		ret = append(ret, r)
	}
	if len(errs) > 0 {
		return nil, fmt.Errorf("roll parsing failed: %v", errs)
	}
	return ret, nil
}

func checkForTrolls(r *RollRequest) string {
	// Validate strings before parsing for weird input.
	switch {
	case r.Multiplier > maxComputedRolls:
		return "I ain't got that many dice."
	case r.Die < 2:
		return fmt.Sprintf("A %d-sided die is pointless, you ass.", r.Die)
	case r.Die > maxDieSize:
		return fmt.Sprintf("A d%d is basically a sphere, wtf.", r.Die)
	case r.Modifier > maxAbsModifier:
		return "You can't add that much to a modifier, that's unreasonable."
	case r.Modifier < -maxAbsModifier:
		return "You can't subtract that much from a modifier, that's unreasonable."
	}
	return ""
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
