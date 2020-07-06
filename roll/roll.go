// Package roll defines a simple request/response interface for requesting dice
// rolls from free-form text and rolling 'dem bones.
package roll

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	rollRegexp = regexp.MustCompile(`(\d*)\s*[dD](\d+)\s*([+-]\s*\d+)?`)

	maxShownRolls    = 10
	maxAbsModifier   = 100
	maxDieSize       = 1000
	maxComputedRolls = 100000000
)

// PlayerRollRequest holds a parsed pre-rolling request. No rolls here.
type PlayerRollRequest struct {
	Timestamp time.Time
	Specs     []*RollSpec
	IsTroll   bool
	TrollMsg  string
}

// PlayerRollResponse holds the post-roll metadata. Lotsa rolls here.
type PlayerRollResponse struct {
	PlayerRollRequest
	Total    int
	Rolls    []RollResult
	IsTroll  bool
	TrollMsg string
}

// RollResult stores metadata about a roll.
type RollResult struct {
	Spec       RollSpec
	Result     int
	BaseRolls  []int
	IsCrit     bool
	IsCritFail bool
}

// RollSpec defines the basic unit of a parsed roll.
type RollSpec struct {
	Multiplier int  // Must be non-negative, zero is a troll.
	Die        int
	Modifier   int
}

type rollTokens struct {
	multiplier string
	die        string
	modifier   string
}

func (pr *PlayerRollResponse) String() string {
	if pr.IsTroll {
		if pr.TrollMsg != "" {
			return pr.TrollMsg
		}
		return "lol, nice try"
	}
	var ss []string
	for _, r := range pr.Rolls {
		ss = append(ss, r.String())
	}
	if len(ss) == 1 {
		return fmt.Sprintf("%s", ss[0])
	} else {
		return fmt.Sprintf("%s Total=**%d**", strings.Join(ss, ", "), pr.Total)
	}
}

func (rr *RollResult) String() string {
	s := []string{rr.Spec.String(), "->"}
	if len(rr.BaseRolls) == 1 && rr.Spec.Modifier == 0 {
		// Format unmodified, single die roll: dXX->Result
		s = append(s, fmt.Sprintf("**%d**", rr.Result))
		return strings.Join(s, "")
	}
	// Format multi-die roll: dXX->r1+r2+...+rn
	s = append(s, "*")
	if len(rr.BaseRolls) == 0 {
		s = append(s, "(nuthin)")
	}
	for i, br := range rr.BaseRolls {
		if i > 0 {
			s = append(s, "+")
		}
		if i >= maxShownRolls {
			s = append(s, fmt.Sprintf("**(%d rolls omitted, ass)**", len(rr.BaseRolls)-i))
			break
		}
		switch {
		case rr.IsCrit:
			s = append(s, fmt.Sprintf("%d(crit)", br))
		case rr.IsCritFail:
			s = append(s, fmt.Sprintf("%d(crit-fail)", br))
		default:
			s = append(s, fmt.Sprintf("%d", br))
		}
	}
	s = append(s, "*")
	// Append modifier
	if rr.Spec.Modifier != 0 {
		s = append(s, fmt.Sprintf("%+d", rr.Spec.Modifier))
	}
	// Append total.
	s = append(s, fmt.Sprintf("=**%d**", rr.Result))
	return strings.Join(s, "")
}

func HasRollRequest(msg string) bool {
	return len(rollRegexp.FindAllStringSubmatch(msg, -1)) > 0
}

func ParseRequest(msg string) (*PlayerRollRequest, error) {
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
		// Validate strings before parsing for weird input.
		switch {
		case len(tk.multiplier) > 20:
			isTroll = true
		case len(tk.die) > 10:
			isTroll = true
			trollMsg = fmt.Sprintf("A d%s is basically a sphere, wtf.", tk.die)
		case len(tk.modifier) > 10:
			isTroll = true
		}

		r, err := tk.parseRollSpec()
		if err != nil {
			errs = append(errs, err)
			continue
		}
		switch {
		case r.Multiplier > maxComputedRolls:
			isTroll = true
			trollMsg = "I ain't got that many dice."
		case r.Multiplier > maxComputedRolls:
			isTroll = true
			trollMsg = "I ain't got that many dice."
		case r.Die == 1:
			isTroll = true
			trollMsg = "A 1-sided die is pointless, you ass."
		case r.Die == 0:
			isTroll = true
			trollMsg = "A 0-sided die is pointless, you ass."
		case r.Die < 0:
			isTroll = true
			trollMsg = "A negative-sided die is pointless, you ass."
		case r.Die > maxDieSize:
			isTroll = true
			trollMsg = fmt.Sprintf("A d%d is basically a sphere, wtf.", r.Die)
		case r.Modifier > maxAbsModifier:
			isTroll = true
			trollMsg = "You can't add that much to a modifier."
		case r.Modifier < -maxAbsModifier:
			isTroll = true
			trollMsg = "You can't subtract that much from a modifier."
		}
		specs = append(specs, r)
	}
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to parse rolls: %v", errs)
	}
	return &PlayerRollRequest{
		Timestamp: time.Now(),
		Specs:     specs,
		IsTroll:   isTroll,
		TrollMsg:  trollMsg,
	}, nil
}

func Roll(req *PlayerRollRequest) PlayerRollResponse {
	ret := PlayerRollResponse{
		PlayerRollRequest: *req,
	}
	if req.IsTroll {
		ret.IsTroll = true
		ret.TrollMsg = req.TrollMsg
		return ret
	}
	for _, s := range req.Specs {
		r := s.Roll()
		ret.Total += r.Result
		ret.Rolls = append(ret.Rolls, r)
	}
	return ret
}

func (rs *RollSpec) Roll() RollResult {
	var sum int
	var br []int
	for i := 0; i < rs.Multiplier; i++ {
		r := rand.Intn(rs.Die) + 1
		br = append(br, r)
		sum += r
	}
	return RollResult{
		Spec:       *rs,
		Result:     sum + rs.Modifier,
		BaseRolls:  br,
		IsCrit:     len(br) == 1 && br[0] == rs.Die,
		IsCritFail: len(br) == 1 && br[0] == 1,
	}
}

func (rs *RollSpec) String() string {
	var tks []string
	if rs.Multiplier != 1 {
		tks = append(tks, strconv.Itoa(rs.Multiplier))
	}
	tks = append(tks, "d"+strconv.Itoa(rs.Die))
	if rs.Modifier != 0 {
		tks = append(tks, fmt.Sprintf("%+d", rs.Modifier))
	}
	return strings.Join(tks, "")
}

func (rt *rollTokens) parseRollSpec() (*RollSpec, error) {
	if rt == nil {
		return nil, fmt.Errorf("nil match")
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
