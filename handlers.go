// Package handlers registers packages for handling Baevents and providing
// Baesponses based on user input.
package dicebae

import (
	"time"

	"dicebae/baepi"
	"dicebae/roll"

	"github.com/bwmarrin/discordgo"
)

func (db *diceBae) initHandlers(args *Baergs) error {
	db.addBaeSaysHandler("roll", roll.NewRollHandler())
	db.addBaeSaysHandler("history", roll.NewHistoryHandler(10, "history"))
	db.addBaeSaysHandler("latest", roll.NewHistoryHandler(1, "latest"))
	return nil
}

func (db *diceBae) addBaeSaysHandler(name string, bh baepi.BaeSayHandler) {
	db.session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			// Don't reply to yourself.
			return
		}
		bf := &baepi.BaestFriend{
			ID:       m.Author.ID,
			Username: m.Author.Username,
		}
		be := &baepi.Baevent{
			Speaker: bf,
			Message: m.Content,
		}
		if !bh.ShouldSay(db, be) {
			// Nothing to say here.
			return
		}
		resp, err := bh.SayWithBae(db, be)
		if err != nil {
			db.LogError("bae can't say! no way: %v", err)
		}
		msg := resp.Message
		if resp.MentionUser {
			msg = bf.Mention(resp.Message)
		}
		s.ChannelMessageSend(m.ChannelID, msg)
		he := &baepi.BaeHistoryEntry{
			HandlerName: name,
			Response:    resp,
			TimeSaid:    time.Now(),
			RepliedTo:   bf,
		}
		db.appendToHistory(he)
		db.LogInfo("Sent response: %#v", resp)
	})
}
