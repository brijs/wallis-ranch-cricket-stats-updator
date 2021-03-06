package main

import (
	"log"
	"sync"
	"flag"
)

func main() {
	log.Printf("Main: Begin WR Stats update...\n")

	// parse flags
	var MAX_PLAYERS int
	flag.IntVar(&MAX_PLAYERS, "max_players", 500, "max number of players to process")
	flag.Parse()

	// print config
	log.Printf("Main: MAX_PLAYERS= %d", MAX_PLAYERS)

	playerStatJobs := make(chan PlayerProfile, 1000)
	outputProfileJobs := make(chan PlayerProfile, 1000)

	// threads to get stats
	const NUM_THREADS = 1
	wg := new(sync.WaitGroup)
	wg.Add(NUM_THREADS)
	for i := 0; i < NUM_THREADS; i++ {
		log.Printf("Main: Starting thread %d to get stats\n", i)
		go GetPlayerStats(playerStatJobs, outputProfileJobs, wg)
	}

	// get list of players
	playerProfiles := GetWRPlayers()
	i := 0
	for _, p := range playerProfiles {
		if i >= MAX_PLAYERS {
			break
		}
		log.Printf("Main: Adding profile %s to get stats\n", p.Name)
		playerStatJobs <- p
		i++
	}
	close(playerStatJobs)

	// wait for
	log.Printf("Main: Waiting for jobs to finish")
	wg.Wait()
	log.Printf("Main: Done Waiting for jobs to finish")
	close(outputProfileJobs)

	// collect output
	var outputProfiles []PlayerProfile
	for o := range outputProfileJobs {
		outputProfiles = append(outputProfiles, o)
		log.Printf("Main: Output => %s", o.Name)
	}

	// Update sheets
	PersistBattingStats(outputProfiles)
	PersistBowlingStats(outputProfiles)
	PersistFieldingStats(outputProfiles)
	log.Printf("Main: Done")
}
