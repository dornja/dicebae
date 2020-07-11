// Package main runs the bae forever.
package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"dicebae"
)

var (
	apiKey       = flag.String("key", "", "The Bot API key, it's a secret to everyone.")
	playerIDList = flag.String("players", "", "A comma-separated list of DNDBeyond player IDs. This is the number in a character sheet URL.")

	maxShownHistory = 10
)

func main() {
	flag.Parse()
	if *apiKey == "" {
		fmt.Println("You need to provide --key=<the dicebae API key>")
		return
	}
	var playerIDs []int
	if ids := *playerIDList; ids != "" {
		for _, tk := range strings.Split(ids, ",") {
			v, err := strconv.ParseInt(tk, 10, 64)
			if err != nil {
				fmt.Printf("Invalid player IDs %q, parse error: %v\n", *playerIDList, err)
				return
			}
			playerIDs = append(playerIDs, int(v))
		}
	}
	db, err := dicebae.NewBae(&dicebae.Baergs{APIKey: *apiKey, PlayerIDs: playerIDs})
	if err != nil {
		fmt.Errorf("Failed to create the bae: %v", err)
	}
	if err := db.LetsRoll(); err != nil {
		fmt.Errorf("This bae won't roll: %v", err)
	}
}
