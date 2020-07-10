package roll

import (
	"sort"
	"strings"

	"dicebae/baepi"
)

type HistoryHandler struct {
	maxEntries int
	hotword    string
}

func NewHistoryHandler(maxEntries int, hotword string) *HistoryHandler {
	return &HistoryHandler{maxEntries: maxEntries, hotword: hotword}
}

func (hh *HistoryHandler) ShouldSay(db baepi.DiceBae, e *baepi.Baevent) bool {
	return strings.HasPrefix(e.Message, "!"+hh.hotword)
}

func (hh *HistoryHandler) SayWithBae(db baepi.DiceBae, e *baepi.Baevent) (*baepi.Baesponse, error) {
	hist := db.FetchHistory(&baepi.BaeHistoKey{HandlerName: "roll"}, 100)
	if len(hist) == 0 {
		return &baepi.Baesponse{Message: "History of what?"}, nil
	}
	// Split by replied-to user.
	histPerBF := make(map[baepi.BaestFriend][]*baepi.BaeHistoryEntry)
	for _, bhe := range hist {
		if bhe.RepliedTo != nil {
			id := *bhe.RepliedTo
			histPerBF[id] = append(histPerBF[id], bhe)
		}
	}
	// Extract unique users and sort to escape nondeterministic map iter order.
	bfs := make([]baepi.BaestFriend, 0, len(histPerBF))
	for u := range histPerBF {
		bfs = append(bfs, u)
	}
	sort.Slice(bfs, func(i, j int) bool {
		return bfs[i].Username < bfs[j].Username
	})

	// Build output string.
	var out []string
	if hh.maxEntries > 1 {
		out = append(out, "**Roll History (newest --> oldest)**")
	} else {
		out = append(out, "**Latest Rolls**")
	}
	for _, bf := range bfs {
		bfHist := histPerBF[bf]
		rs := []string{}
		for i, bhe := range bfHist {
			if i >= hh.maxEntries {
				break
			}
			rs = append(rs, bhe.Response.Message)
		}
		out = append(out, bf.Mention("["+strings.Join(rs, "] [")+"]"))
	}
	return &baepi.Baesponse{
		Message: strings.Join(out, "\n"),
	}, nil
}
