// Package roll defines a simple request/response interface for requesting dice
// rolls from free-form text and rolling 'dem bones.
package roll

import (
	"math/rand"
	"regexp"
	"strings"
	"time"

	"dicebae/baepi"
)

var (
	maxResponseLength = 10
	maxShownRolls     = 10
	maxDieSize        = 1000
	maxAbsModifier    = 10000
	maxComputedRolls  = 1000000
	rollRegexp        = regexp.MustCompile(`(\d*)\s*[dD](\d+)\s*([+-]\s*\d+)?`)
)

// RollHandler implements the BaeSayHandler interface for rolling 'dem bones.
type RollHandler struct {
	kelgwynFrustrator *rand.Rand
}

// RollRequest stores a parsed user roll request (XdN[+|-]Mod) along with a
// troll message if the request was dumb.
type RollRequest struct {
	Multiplier int
	Die        int
	Modifier   int
	TrollMsg   string
}

// RollResult stores the outcome of rolling a single RollRequest.
type RollResult struct {
	Request    *RollRequest
	Result     int
	BaseRolls  []int
	IsCrit     bool
	IsCritFail bool
}

// RollResult stores the outcome of rolling potentially many RollRequest, and
// represents a single response to a user's request to roll one to many
// RollRequests. If any request was dumb, the response is a troll response.
type RollResponse struct {
	Total         int
	Results       []*RollResult
	TrollResponse string
}

func NewRollHandler() *RollHandler {
	return &RollHandler{
		kelgwynFrustrator: rand.New(rand.NewSource(int64(time.Now().Nanosecond()))),
	}
}

func (rh *RollHandler) ShouldSay(db baepi.DiceBae, e *baepi.Baevent) bool {
	return len(rollRegexp.FindAllStringSubmatch(e.Message, -1)) > 0
}

func (rh *RollHandler) SayWithBae(db baepi.DiceBae, e *baepi.Baevent) (*baepi.Baesponse, error) {
	reqs, err := parseRollRequests(e.Message)
	if err != nil {
		return nil, err
	}

	// Roll 'dem bones.
	var resp RollResponse
	var trolls []string
	for _, req := range reqs {
		res := req.Roll(rh.kelgwynFrustrator)
		resp.Total += res.Result
		resp.Results = append(resp.Results, res)
		if req.TrollMsg != "" {
			trolls = append(trolls, req.TrollMsg)
		}
	}
	switch {
	case len(reqs) > maxResponseLength:
		resp.TrollResponse = "I refuse to do that much work, ass."
	case len(trolls) > 0:
		resp.TrollResponse = strings.Join(trolls, " Also: ")
	}

	return &baepi.Baesponse{
		Message:         resp.String(),
		MentionUser:     true,
		HandlerMetadata: resp,
	}, nil
}

func (rs *RollRequest) Roll(rng *rand.Rand) *RollResult {
	if rs.TrollMsg != "" {
		return &RollResult{
			Request:    rs,
			Result:     1,
			IsCritFail: true,
		}
	}
	var sum int
	var br []int
	for i := 0; i < rs.Multiplier; i++ {
		r := rng.Intn(rs.Die) + 1
		br = append(br, r)
		sum += r
	}
	return &RollResult{
		Request:    rs,
		Result:     sum + rs.Modifier,
		BaseRolls:  br,
		IsCrit:     len(br) == 1 && br[0] == rs.Die,
		IsCritFail: len(br) == 1 && br[0] == 1,
	}
}
