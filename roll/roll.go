// Package roll defines a simple request/response interface for requesting dice
// rolls from free-form text and rolling 'dem bones.
package roll

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var (
	maxShownRolls = 10
)

// RollRequest holds a parsed pre-rolling request. No rolls here.
type RollRequest struct {
	Timestamp time.Time
	Specs     []*RollSpec
	IsTroll   bool
	TrollMsg  string
}

// RollResponse holds the post-roll metadata. Lotsa rolls here.
type RollResponse struct {
	RollRequest
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
	Multiplier int // Must be non-negative, zero is a troll.
	Die        int
	Modifier   int
}

func (pr *RollResponse) String() string {
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

func Roll(req *RollRequest) RollResponse {
	ret := RollResponse{
		RollRequest: *req,
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
