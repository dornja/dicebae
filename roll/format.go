package roll

import (
	"fmt"
	"strconv"
	"strings"
)

func (rs *RollRequest) String() string {
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

func (rr *RollResult) String() string {
	s := []string{rr.Request.String(), "->"}
	if len(rr.BaseRolls) == 1 && rr.Request.Modifier == 0 {
		// Format unmodified, single die roll: dXX->Result
		switch {
		case rr.IsCrit:
			s = append(s, fmt.Sprintf("**%d (Crit!)**", rr.Result))
		case rr.IsCritFail:
			s = append(s, fmt.Sprintf("**%d (Crit-Fail!)**", rr.Result))
		default:
			s = append(s, fmt.Sprintf("**%d**", rr.Result))
		}
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
	if rr.Request.Modifier != 0 {
		s = append(s, fmt.Sprintf("%+d", rr.Request.Modifier))
	}
	// Append total.
	s = append(s, fmt.Sprintf("=**%d**", rr.Result))
	return strings.Join(s, "")
}

func (rr *RollResponse) String() string {
	if rr.TrollResponse != "" {
		return rr.TrollResponse
	}
	var ss []string
	for _, r := range rr.Results {
		ss = append(ss, r.String())
	}
	if len(ss) == 1 {
		return fmt.Sprintf("%s", ss[0])
	} else {
		return fmt.Sprintf("%s Total=**%d**", strings.Join(ss, ", "), rr.Total)
	}
}
