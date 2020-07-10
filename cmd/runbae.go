// Package main runs the bae forever.
package main

import (
	"flag"
	"fmt"

	"dicebae"
)

var (
	apiKey = flag.String("key", "", "The Bot API key, it's a secret to everyone.")

	maxShownHistory = 10
)

func main() {
	flag.Parse()
	if *apiKey == "" {
		fmt.Println("You need to provide --key=<the dicebae API key>")
		return
	}
	db, err := dicebae.NewBae(&dicebae.Baergs{APIKey: *apiKey})
	if err != nil {
		fmt.Errorf("Failed to create the bae: %v", err)
	}
	if err := db.LetsRoll(); err != nil {
		fmt.Errorf("This bae won't roll: %v", err)
	}
}

/*
func handleRollRequest(s *discordgo.Session, m *discordgo.MessageCreate) {
	p := Player{
		ID:   m.Author.ID,
		Name: m.Author.Username,
	}
	msg := m.Content
	req, err := roll.ParseRollRequest(msg)
	if err != nil {
		msg := fmt.Sprintf("%s: I beefed it on your roll: %v", p.AtPlayer(), err)
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}
	r := roll.Roll(req)

	if !r.IsTroll {
		history.ResponsesByID[p] = append(history.ResponsesByID[p], r)
	}
	fmt.Printf("%s: %s\n", p.AtPlayer(), r.String())
	fmtSay(s, m, fmt.Sprintf("%s: %s", p.AtPlayer(), r.String()))
}

func handleShowLatest(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(history.ResponsesByID) == 0 {
		fmtSay(s, m, "Latest what?")
		return
	}
	out := []string{"**Latest Rolls Per Player**"}
	for p, roll := range history.ResponsesByID {
		if len(roll) > 0 {
			r := roll[len(roll)-1]
			out = append(out, fmt.Sprintf("%s %s", p.AtPlayer(), r.String()))
		}
	}
	fmtSay(s, m, strings.Join(out, "\n"))
}

func handleShowHistory(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(history.ResponsesByID) == 0 {
		fmtSay(s, m, "History of what?")
		return
	}
	out := []string{"**Player Roll History (newest first)**"}
	for p, resps := range history.ResponsesByID {
		rs := []string{p.AtPlayer()}
		for i := range resps {
			if i >= maxShownHistory {
				break
			}
			ri := len(resps) - 1 - i
			rs = append(rs, resps[ri].String())
		}
		out = append(out, strings.Join(rs, ", "))
	}
	fmtSay(s, m, strings.Join(out, "\n"))
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	isAtBot := strings.Contains(m.Content, s.State.User.ID)
	isCommand := isAtBot || strings.HasPrefix(m.Content, "!")

	switch {
	case roll.HasRollRequest(m.Content):
		handleRollRequest(s, m)
	case isCommand && strings.Contains(strings.ToLower(m.Content), "hist"):
		handleShowHistory(s, m)
	case isCommand && strings.Contains(strings.ToLower(m.Content), "latest"):
		handleShowLatest(s, m)
	case m.Content == "ping":
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	case m.Content == "ping":
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}*/
