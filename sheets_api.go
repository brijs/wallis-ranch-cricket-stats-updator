package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func getProfileURL(p PlayerProfile) string {
	// https://cricheroes.in/player-profile/1140525/Brijesh-Shetty?
	return fmt.Sprintf("https://cricheroes.in/player-profile/%d/%s?", p.ID, strings.Replace(p.Name, " ", "-", -1))
}

func PersistBattingStats(profiles []PlayerProfile) {
	rangeData := "Batting!A1:Z1000" // 1000 rows
	values := [][]interface{}{
		{"", "Name", "Matches", "Innings", "Runs", "Not Outs", "Highest",
			"Avg", "SR", "4s", "6s", "0s", "30s", "60s", "100s"}}
	for _, p := range profiles {
		values = append(values, []interface{}{fmt.Sprintf("=IMAGE(\"%s\")", p.PhotoURL),
			fmt.Sprintf("=HYPERLINK(\"%s\",\"%s\")", getProfileURL(p), p.Name), p.Batting.Matches, p.Batting.Innings, p.Batting.Runs, p.Batting.NotOuts, p.Batting.Highest,
			p.Batting.Average, p.Batting.StrikeRate, p.Batting.Fours, p.Batting.Sixes, p.Batting.Ducks, p.Batting.Thirties, p.Batting.Fifties, p.Batting.Hundreds})
	}

	persistToSheet(rangeData, values)
	log.Println("Sheets: Done updating Batting Stats.")
}

func PersistBowlingStats(profiles []PlayerProfile) {
	rangeData := "Bowling!A1:Z1000" // 1000 rows
	values := [][]interface{}{
		{"", "Name", "Matches", "Innings", "Overs", "Maidens", "Runs",
			"Wickets", "Economy", "SR", "Avg", "Wides", "NoBalls", "Dots", "Best"}}
	for _, p := range profiles {
		values = append(values, []interface{}{fmt.Sprintf("=IMAGE(\"%s\")", p.PhotoURL),
			fmt.Sprintf("=HYPERLINK(\"%s\",\"%s\")", getProfileURL(p), p.Name), p.Bowling.Matches, p.Bowling.Innings, p.Bowling.Overs, p.Bowling.Maidens, p.Bowling.Runs,
			p.Bowling.Wickets, p.Bowling.Economy, p.Bowling.StrikeRate, p.Bowling.Average, p.Bowling.Wides, p.Bowling.NoBalls, p.Bowling.DotBalls, p.Bowling.BestBowling})
	}

	persistToSheet(rangeData, values)
	log.Println("Sheets: Done updating Bowling Stats.")
}

func PersistFieldingStats(profiles []PlayerProfile) {
	rangeData := "Fielding!A1:Z1000" // 1000 rows
	values := [][]interface{}{
		{"", "Name", "Matches", "Catches", "Ct Behind", "Stumpings", "RunOuts", "AssistedRunOuts"}}
	for _, p := range profiles {
		values = append(values, []interface{}{fmt.Sprintf("=IMAGE(\"%s\")", p.PhotoURL),
			fmt.Sprintf("=HYPERLINK(\"%s\",\"%s\")", getProfileURL(p), p.Name), p.Fielding.Matches, p.Fielding.Catches, p.Fielding.CaughtBehind, p.Fielding.Stumpings,
			p.Fielding.RunOuts, p.Fielding.AssistedRunOuts})
	}

	persistToSheet(rangeData, values)
	log.Println("Sheets: Done updating Bowling Stats.")
}

func persistToSheet(rangeData string, values [][]interface{}) {
	ctx := context.Background()

	b, err := ioutil.ReadFile("google_sheets_credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// only scoped to edit files created by WR Cricket client
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/drive.file")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	spreadsheetId2 := "1nE9NtBm-95XVp1XptPgLpSsBEyialbcvnJN_2f1PIeg"
	rb := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: "USER_ENTERED",
	}
	rb.Data = append(rb.Data, &sheets.ValueRange{
		Range:  rangeData,
		Values: values,
	})
	_, err = srv.Spreadsheets.Values.BatchUpdate(spreadsheetId2, rb).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}
	return
}
