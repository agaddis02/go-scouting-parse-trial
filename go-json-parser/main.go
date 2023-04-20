package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/imdario/mergo"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
	//	"time"
)

type Action struct {
	X         float64 `json:"x"`
	Y         float64 `json:"Y"`
	Timestamp int     `json:"time_stamp"`
	Type      string  `json:"type"`
}

type Actions struct {
	ListOfActions []Action `json:"actions"`
}
type Scoring struct {
	Teleop     Actions        `json:"teleop"`
	Autonomous Actions        `json:"autonomous"`
	PostMatch  map[string]any `json:"post_match"`
}
type MatchJson struct {
	MatchNumber    int     `json:"match_number"`
	TeamNumber     int     `json:"team_number"`
	ConfiguredTeam string  `json:"configured_team"`
	EventKey       string  `json:"event_key"`
	MatchKey       string  `json:"match_key"`
	ScouterId      string  `json:"scouter_id"`
	ScoringSilos   Scoring `json:"scoring"`
}

type SummaryTable struct {
	MatchNumber    int
	TeamNumber     int
	EventKey       string
	Year           int
	ConfiguredTeam string
	MatchKey       string
	ScouterId      string
	Columns        map[string]int
	PostMatch      map[string]any
}

const (
	FILENAME = "matches/3/1153_red_1_WilliamGohi.json"
)

func flattenScoringActions(actions []Action, columnPrefix string) SummaryTable {
	var summaryTable SummaryTable
	summaryTable.Columns = make(map[string]int)
	for i := 0; i < len(actions); i++ {
		summaryTable.Columns[columnPrefix+actions[i].Type] += 1
	}
	return summaryTable
}

func main() {

	err := godotenv.Load()

	host := os.Getenv("HOST")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}
	user := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DATABASE")

	if err != nil {
		log.Fatal("Error loading .env file")
	}


	// open json file
	var matchData MatchJson
	fileBytes, _ := os.ReadFile(FILENAME)
	err = json.Unmarshal(fileBytes, &matchData)

	if err != nil {
		log.Fatal(err)
	}

	teleopActions := matchData.ScoringSilos.Teleop.ListOfActions
	autonActions := matchData.ScoringSilos.Autonomous.ListOfActions
	postMatchActions := matchData.ScoringSilos.PostMatch

	summaryTelop := flattenScoringActions(teleopActions, "TELEOP.")
	summaryAutonomous := flattenScoringActions(autonActions, "AUTONOMOUS.")

	var tableToUpload SummaryTable
	tableToUpload.PostMatch = postMatchActions
	tableToUpload.Year = 2024
	tableToUpload.MatchNumber = matchData.MatchNumber
	tableToUpload.ConfiguredTeam = matchData.ConfiguredTeam
	tableToUpload.EventKey = matchData.EventKey
	tableToUpload.ScouterId = matchData.ScouterId
	tableToUpload.MatchKey = matchData.MatchKey

	err = mergo.Merge(&tableToUpload, summaryTelop)

	if err != nil {
		log.Fatal(err)
	}

	err = mergo.Merge(&tableToUpload, summaryAutonomous)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(tableToUpload)


	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	result, err := db.Exec("SELECT * FROM public.auton_actions LIMIT 5")

	// db.QueryRowContext() - Run a SELECT statement, to see if the id of the record is in the table or not
	// if the record is not in the table, we will insert, otherwise we will update
//	db.QueryRowContext()



	db.QueryRowContext()
	err = db.Close()
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(result.RowsAffected())

}
