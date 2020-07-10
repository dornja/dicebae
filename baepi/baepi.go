// Package baepi defines a simple api for interactions between the bae
// and any handlers that define bae features.
package baepi

import (
	"fmt"
	"time"
)

// DiceBae defines the public interface for the bae. This is what handlers will
// interact with instead of the raw, powerful bae implementation.
type DiceBae interface {
	LetsRoll() error
	LogInfo(string, ...interface{})
	LogError(string, ...interface{})
	FetchHistory(*BaeHistoKey, int) []*BaeHistoryEntry
}

// BaeSayHandler defines the interface for a simple handler that conditionally
// responds to anyone in a channel containing the bae. Bae handlers need to
// only define under what conditions to trigger, and what to say if triggered.
// Mostly, this abstracts away the extra, unnecessary details provided by
// discordgo, since most of the time all the bae wants to do is reply to some
// hotword or regexp spoken by a player.
type BaeSayHandler interface {
	// ShouldSay returns whether a handler should respond. If this returns true,
	// SayWithBae will be called.
	ShouldSay(DiceBae, *Baevent) bool
	// SayWithBae returns what the bae should say when ShouldSay returns OK.
	SayWithBae(DiceBae, *Baevent) (*Baesponse, error)
}

// Baevent is a Bae Event. Specifically, it encapsulates a user sending a
// message in a discord channel containing the bae.
type Baevent struct {
	Speaker *BaestFriend
	Message string
}

// Baesponse contains the bae's response to a Baevent. Beyond the message to be
// sent to the channel, it also includes a generic metadata argument that will
// be stored in the bae's history. Handlers can use this metadata by pulling
// old replies out of the bae's history.
type Baesponse struct {
	Message         string
	MentionUser     bool
	HandlerMetadata interface{}
}

// BaestFriend defines a user entity in discord. The ID can be used to <@ID>
// mention a user in a Baesponse and the username is the human-readable
// username. This is essentially a subset of the User fields from discordgo.
type BaestFriend struct {
	ID       string
	Username string
}

// BaeHistoKey contains optional constraints when searching the bae's history.
type BaeHistoKey struct {
	BaestFriendID string
	HandlerName   string
}

// BaeHistoryEntry contains a single bae response along with metadata about what handler
type BaeHistoryEntry struct {
	HandlerName string
	Response    *Baesponse
	TimeSaid    time.Time
	RepliedTo   *BaestFriend
}

// Mention returns a modified message string that will trigger a mention, e.g.,
// @someuser: message, in Discord.
func (bf *BaestFriend) Mention(message string) string {
	return fmt.Sprintf("<@%s> %s", bf.ID, message)
}

func (he *BaeHistoryEntry) Matches(k *BaeHistoKey) bool {
	switch {
	case k.HandlerName != "" && k.HandlerName != he.HandlerName:
		return false
	case k.BaestFriendID != "" && k.BaestFriendID != he.RepliedTo.ID:
		return false
	default:
		return true
	}
}
