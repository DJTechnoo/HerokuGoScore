package main

import(

	"fmt"
	"net/http"
	"os"
	"log"
	"strconv"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"encoding/json"

)

// Database constants
const dbName 		= "highscores"
const dbCollection 	= "scores"
const dbURL			= "mongodb://scoreuser:spillprog4life@ds145083.mlab.com:45083/highscores"


type scoreEntry struct {
	ID bson.ObjectId 	`bson:"_id,omitempty" json:"-"`
	Score int 			`bson:"score" json:"score"`
}


type scoreResponse struct {
	HighScore int		`json:"high_score"`
	OtherScores []int	`json:"other_scores"` 
}



//	Adds scores to the DB as entries
func addToDB(entry scoreEntry) {

	session, err := mgo.Dial(dbURL)
	if err != nil {
		return
	}

	defer session.Close()

	err = session.DB(dbName).C(dbCollection).Insert(entry)

	if err != nil {
		return
	}
}


func createResponse(w http.ResponseWriter, sortedScores []int){
	response := scoreResponse{}				// 1 empty data
	response.HighScore = sortedScores[0]	// take the highest score
	response.OtherScores = sortedScores[1 : len(sortedScores)]	// add rest
	
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	enc.Encode(&response)
}


func sortByScore(w http.ResponseWriter){
	session, err := mgo.Dial(dbURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()


	c := session.DB(dbName).C(dbCollection)

	item := scoreEntry{}

	find := c.Find(nil).Sort("-score")

	sortedScores := make([]int, 0)
	items := find.Iter()
	for items.Next(&item) {
		sortedScores = append(sortedScores, item.Score)
	}


	createResponse(w, sortedScores)
}


//	handles the score
func processScore(score int){
	entry := scoreEntry{ID: bson.NewObjectId(), Score: score}	
						
	addToDB(entry)
}


func scoreHandler(w http.ResponseWriter, r * http.Request){

	switch r.Method {
	case http.MethodGet:
		sortByScore(w)
		
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			return
		}

		latestScore, _ := strconv.Atoi(r.FormValue("score"))
		fmt.Fprintln(w, latestScore)
		processScore(latestScore)
	}
}



func main(){
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	
	
	http.HandleFunc("/scores", scoreHandler)
	http.ListenAndServe(":" + port, nil)

}
