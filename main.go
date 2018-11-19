package main

import(

	//"fmt"
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
	Username string     `bson:"username" json:"username"`
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





func sortByScore(w http.ResponseWriter){
	session, err := mgo.Dial(dbURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()


	c := session.DB(dbName).C(dbCollection)

	sortedScores := []scoreEntry{}

	c.Find(bson.M{}).Sort("-score").Limit(5).All(&sortedScores)

	http.Header.Add(w.Header(), "content-type", "application/json")
	json.NewEncoder(w).Encode(&sortedScores)
	
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
		username   := r.FormValue("username")

		s := scoreEntry{bson.NewObjectId(), latestScore, username}
		
		addToDB(s)
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
