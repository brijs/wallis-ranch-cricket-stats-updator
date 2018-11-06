package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "time"
    "sync"
    "os"
)

type PlayerID uint

type PlayerProfile struct {
    ID       PlayerID
    Name     string
    PhotoURL string
    Batting BattingStat
    Bowling BowlingStat
}

type ApiResponseData struct {
    Players []PlayerProfile
    Batting []BattingStat
    Bowling []BowlingStat
}

type BattingStat struct {
    Matches    uint
    Innings    uint
    NotOuts    uint
    Runs       uint
    Highest    uint
    Average    float32
    StrikeRate float32
    Fours      uint
    Sixes      uint
    Ducks      uint
    Thirties   uint
    Fifties    uint
    Hundreds   uint
}

type BowlingStat struct {
    Matches uint
    Innings uint
    Overs float32
    Runs uint
    Maidens uint
    Wickets uint
    Economy float32
    StrikeRate float32
    Average float32
    Wides uint
    NoBalls uint
    DotBalls uint
    BestBowling string

}

type ApiResponse struct {
    Success bool
    Data    ApiResponseData
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getWrapAPIKey() (string) {
    val, ok := os.LookupEnv("BRIJ_WRAP_API_KEY")
    if !ok {
        log.Fatalln("BRIJ_WRAP_API_KEY is not set")
        return val
    } else {
        return val
    }
}

func GetWRPlayers() ([]PlayerProfile){
    WRPlayersURL := fmt.Sprintf("https://wrapapi.com/use/brij/tests/WRCricPlayers/0.0.4?wrapAPIKey=%s", getWrapAPIKey())
    r, err := myClient.Get(WRPlayersURL)
    defer r.Body.Close()
    if err != nil {
        log.Fatalln(err)
    }

    if r.StatusCode != http.StatusOK {
        log.Fatalln("Received Bad HTTP Status:", r.StatusCode)
    }

    bodyBytes, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Fatalln(err)
    }

    resp := &ApiResponse{}
    err = json.Unmarshal(bodyBytes, resp)
    if err != nil {
        fmt.Println(err)
    }
    if resp.Success != true {
        log.Fatalf("Got Response.success=false %+v", resp)
    }
    if resp.Data.Players == nil {
        log.Fatalf("Got Nil/Empty players %+v", resp)
    }

    return resp.Data.Players

}

func GetPlayerStats(playerStatJobs <-chan PlayerProfile, outputProfiles chan<- PlayerProfile, wg *sync.WaitGroup) {
        for p := range playerStatJobs {
            log.Printf("GetPlayerStats: processing ID=%s\n", p.Name)
            WRPlayerStatsURL := fmt.Sprintf("https://wrapapi.com/use/brij/tests/WRCricPlayerStat/latest?player_id=%d&wrapAPIKey=%s", p.ID, getWrapAPIKey())
            r, err := myClient.Get(WRPlayerStatsURL)
            defer r.Body.Close()
            if err != nil {
                log.Fatalln(err)
            }

            if r.StatusCode != http.StatusOK {
                log.Fatalln("Received Bad HTTP Status:", r.StatusCode)
            }

            bodyBytes, err := ioutil.ReadAll(r.Body)
            if err != nil {
                log.Fatalln(err)
            }

            resp := &ApiResponse{}
            err = json.Unmarshal(bodyBytes, resp)
            if err != nil {
                fmt.Println(err)
            }
            // if resp.Success != true {
            //     log.Fatalf("Got Response.success=false %+v for ID=%d Name=%s", resp, p.ID, p.Name)
            // }
            // if resp.Data.Batting == nil || resp.Data.Bowling == nil {
            //     log.Fatalf("Got Nil Batting/Bowling stat %+v", resp)
            // }

            if len(resp.Data.Batting) > 0 {
                p.Batting = resp.Data.Batting[0]
            }
            if len(resp.Data.Bowling) > 0 {
                p.Bowling = resp.Data.Bowling[0]
            }

            // log.Printf("GetPlayerStats: Adding ID=%+v\n", p)
            outputProfiles <- p
        }

        wg.Done()
        return

}

