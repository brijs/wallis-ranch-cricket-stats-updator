package main

import (
	"math"
)

func GetOldPlayerProfiles() []PlayerProfile {
	var profiles = []PlayerProfile{
		// Siva's old Profile
		{
			ID:       1140896,
			Name:     "Siva (old)",
			PhotoURL: "https://media.cricheroes.in/default/user_profile.png",
		},
	}

	return profiles
}

func max(x, y uint) uint {
	if x > y {
		return x
	}
	return y
}

func numBalls(overs float32) uint {
	var o, b = math.Modf(float64(overs))
	return uint(o*6 + math.Round(b*10))
}

func roundTo2Dec(n float32) float32 {
	return float32(math.Round(float64(n*100)) / 100)
}

func MergeOldStatsForSiva(newProfile PlayerProfile) PlayerProfile {
	// Siva's old stats
	var oldProfile = &PlayerProfile{
		Batting: BattingStat{
			Matches:    16,
			Innings:    16,
			NotOuts:    0,
			Runs:       101,
			Highest:    20,
			Fours:      4,
			StrikeRate: 68.71,
		},
		Bowling: BowlingStat{
			Matches:  16,
			Innings:  16,
			Overs:    41,
			Maidens:  1,
			Runs:     166,
			Wickets:  14,
			Wides:    27,
			NoBalls:  11,
			DotBalls: 155,
		},
		Fielding: FieldingStat{
			Matches: 16,
			Catches: 11,
		},
	}

	var mergedProfile = newProfile

	// BattingStat
	var mBat, nBat, oBat = &mergedProfile.Batting, &newProfile.Batting, &oldProfile.Batting
	mBat.Matches += oBat.Matches
	mBat.Innings += oBat.Innings
	mBat.NotOuts += oBat.NotOuts
	mBat.Runs += oBat.Runs
	mBat.Highest = max(mBat.Highest, oBat.Highest)
	mBat.Fours += oBat.Fours
	mBat.Sixes += oBat.Sixes
	mBat.Thirties += oBat.Thirties
	mBat.Fifties += oBat.Fifties
	mBat.Hundreds += oBat.Hundreds
	mBat.Ducks += oBat.Ducks

	mBat.Average = roundTo2Dec(float32(mBat.Runs) / float32(mBat.Innings-mBat.NotOuts))
	mBat.StrikeRate = roundTo2Dec(float32(mBat.Runs) / ((float32(oBat.Runs) / oBat.StrikeRate) + (float32(nBat.Runs) / nBat.StrikeRate)))

	// BowlingStat
	var mBowl, oBowl = &mergedProfile.Bowling, &oldProfile.Bowling
	mBowl.Matches += oBowl.Matches
	mBowl.Innings += oBowl.Innings
	mBowl.Overs += oBowl.Overs
	mBowl.Runs += oBowl.Runs
	mBowl.Maidens += oBowl.Maidens
	mBowl.Wickets += oBowl.Wickets
	mBowl.Wides += oBowl.Wides
	mBowl.NoBalls += oBowl.NoBalls
	mBowl.DotBalls += oBowl.DotBalls

	mBowl.Economy = roundTo2Dec(float32(mBowl.Runs*6.0) / float32(numBalls(mBowl.Overs)))
	mBowl.StrikeRate = roundTo2Dec(float32(numBalls(mBowl.Overs)) / float32(mBowl.Wickets))
	mBowl.Average = roundTo2Dec(float32(mBowl.Runs) / float32(mBowl.Wickets))

	// FieldingStat
	mergedProfile.Fielding.Catches += oldProfile.Fielding.Catches

	return mergedProfile
}
