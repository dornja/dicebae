// History handles storing previous Baesponses in-memory and exposes functions
// for handlers to query earlier responses.
package dicebae

import (
	"dicebae/baepi"
)

var (
	maxHistoryEntries = 1000
)

// FetchHistory returns up to n BaeHistoryEntries matching the constraints in
// the given BaeHistoKey, from most recent to oldest. All fields in the key are
// considered optional.
func (db *diceBae) FetchHistory(k *baepi.BaeHistoKey, n int) []*baepi.BaeHistoryEntry {
	if n > maxHistoryEntries {
		n = maxHistoryEntries
	}
	// We could be more clever about this and actually index the history, but a
	// linear scan should be fast enough given the limit on history size.
	var ret []*baepi.BaeHistoryEntry
	for i := range db.history {
		if e := db.history[len(db.history)-1-i]; e.Matches(k) {
			ret = append(ret, e)
		}
		if len(ret) >= n {
			db.LogInfo("Terminating history search early.")
			break
		}
	}
	return ret
}

func (db *diceBae) appendToHistory(he *baepi.BaeHistoryEntry) {
	// This implementation is currently very dumb, but is broken out here to make
	// adding things like a per-channel or per-user index easier in the future.
	db.history = append(db.history, he)
}
