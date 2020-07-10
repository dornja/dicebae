// Package dicebae implements the baepi DiceBae interface. It essentially wraps
// a discordgo session to provide a more bae-focused experience. Specifically,
// it defines a bot centered around responding to hotwords spoken in discord
// channels with simple, easy-to-write handlers.
package dicebae

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"dicebae/baepi"

	"github.com/bwmarrin/discordgo"
)

// Baergs contains arguments for the creation of the bae.
type Baergs struct {
	APIKey string // Required.
	LogDir string
}

// diceBae implements the DiceBae interface defined in the baepi.
type diceBae struct {
	session *discordgo.Session
	logFile *os.File
	logger  *log.Logger
	history []*baepi.BaeHistoryEntry
}

// NewBae returns a hot, fresh bae with validated and initialized handlers.
func NewBae(args *Baergs) (baepi.DiceBae, error) {
	dg, err := discordgo.New("Bot " + args.APIKey)
	if err != nil {
		return nil, fmt.Errorf("Error: %v", err)
	}

	db := &diceBae{
		session: dg,
		history: make([]*baepi.BaeHistoryEntry, 0, 1024),
	}
	if err := db.initHandlers(args); err != nil {
		return nil, fmt.Errorf("bae failed to init handlers: %v", err)
	}
	if err := db.initLogger(args.LogDir); err != nil {
		return nil, fmt.Errorf("bae failed to init logger: %v", err)
	}
	return db, nil
}

// LetsRoll initializes a discordgo session and will attempt to serve responses
// from each dicebae handler until killed.
func (db *diceBae) LetsRoll() error {
	if err := db.session.Open(); err != nil {
		return fmt.Errorf("failed to open Discord session: %v", err)
	}
	defer db.session.Close()
	defer db.logFile.Close()

	db.LogInfo("I have no dice, but I must roll. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	// Block on this channel until we get a termination signal.
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	db.LogInfo("Later dopes.")
	return nil
}
