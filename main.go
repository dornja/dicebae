package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"dicebae/roll"

	"github.com/bwmarrin/discordgo"
)

var (
	apiKey = flag.String("key", "", "The Bot API key, it's a secret to everyone.")

	history         RollHistory
	maxShownHistory = 10
)

// RollHistory should actually be a singleton or something, it stores rolls. It
// should also have a goroutine that syncs it out to disk now and then, and
// restores state upon at start. Also we're not using PlayersSeen, it should be
// a deterministic, sorted order for the map. Also, could build some cool stuff
// to track NPCs rolling.
type RollHistory struct {
	PlayersSeen   []*Player
	ResponsesByID map[Player][]roll.PlayerRollResponse
}

type Player struct {
	ID   string
	Name string
}

func main() {
	flag.Parse()
	if *apiKey == "" {
		fmt.Println("You need to provide --key=<the dicebae API key>")
		return
	}
	history.ResponsesByID = make(map[Player][]roll.PlayerRollResponse)
	dg, err := discordgo.New("Bot " + *apiKey)
	if err != nil {
		fmt.Errorf("Error: %v", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	// TODO: automatically retry failed connections a few times with a backoff.
	fmt.Println("I have no dice, but I must roll. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func fmtSay(s *discordgo.Session, m *discordgo.MessageCreate, msg string, args ...interface{}) {
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(msg, args...))
}

func handleRollRequest(s *discordgo.Session, m *discordgo.MessageCreate) {
	p := Player{
		ID:   m.Author.ID,
		Name: m.Author.Username,
	}
	msg := m.Content
	req, err := roll.ParseRequest(msg)
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

func (p Player) AtPlayer() string {
	return fmt.Sprintf("<@%s>", p.ID)
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
}
