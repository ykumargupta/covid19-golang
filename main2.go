package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"io/ioutil"
)

var client *mongo.Client

type Meta struct {
	Date         string `json:"date,omitempty" bson:"date,omitempty"`
	Last_updated string `json:"last_updated,omitempty" bson:"last_updated,omitempty"`
}
type Total struct {
	Confirmed   int `json:"confirmed,omitempty" bson:"confirmed,omitempty"`
	Deceased    int `json:"deceased,omitempty" bson:"deceased,omitempty"`
	Recovered   int `json:"recovered,omitempty" bson:"recovered,omitempty"`
	Tested      int `json:"tested,omitempty" bson:"tested,omitempty"`
	Vaccinated1 int `json:"vaccinated1,omitempty" bson:"vaccinated1,omitempty"`
	Vaccinated2 int `json:"vaccinated2,omitempty" bson:"vaccinated2,omitempty"`
}

type State struct {
	Meta  Meta  `json:"meta,omitempty" bson:"meta,omitempty"`
	Total Total `json:"total,omitempty" bson:"total,omitempty"`
}

type India struct {
	// ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ID string `json:"id"`
	AN State  `json:"AN,omitempty" bson:"AN,omitempty"`
	AP State  `json:"AP,omitempty" bson:"AP,omitempty"`
	AR State  `json:"AR,omitempty" bson:"AR,omitempty"`
	AS State  `json:"AS,omitempty" bson:"AS,omitempty"`
	BR State  `json:"BR,omitempty" bson:"BR,omitempty"`
	CH State  `json:"CH,omitempty" bson:"CH,omitempty"`
	CT State  `json:"CT,omitempty" bson:"CT,omitempty"`
	DL State  `json:"DL,omitempty" bson:"DL,omitempty"`
	DN State  `json:"DN,omitempty" bson:"DN,omitempty"`
	GA State  `json:"GA,omitempty" bson:"GA,omitempty"`
	GJ State  `json:"GJ,omitempty" bson:"GJ,omitempty"`
	HP State  `json:"HP,omitempty" bson:"HP,omitempty"`
	HR State  `json:"HR,omitempty" bson:"HR,omitempty"`
	JH State  `json:"JH,omitempty" bson:"JH,omitempty"`
	JK State  `json:"JK,omitempty" bson:"JK,omitempty"`
	KA State  `json:"KA,omitempty" bson:"KA,omitempty"`
	KL State  `json:"KL,omitempty" bson:"KL,omitempty"`
	LA State  `json:"LA,omitempty" bson:"LA,omitempty"`
	LD State  `json:"LD,omitempty" bson:"LD,omitempty"`
	MH State  `json:"MH,omitempty" bson:"MH,omitempty"`
	ML State  `json:"ML,omitempty" bson:"ML,omitempty"`
	MN State  `json:"MN,omitempty" bson:"MN,omitempty"`
	MP State  `json:"MP,omitempty" bson:"MP,omitempty"`
	MZ State  `json:"MZ,omitempty" bson:"MZ,omitempty"`
	NL State  `json:"NL,omitempty" bson:"NL,omitempty"`
	OR State  `json:"OR,omitempty" bson:"OR,omitempty"`
	PB State  `json:"PB,omitempty" bson:"PB,omitempty"`
	PY State  `json:"PY,omitempty" bson:"PY,omitempty"`
	RJ State  `json:"RJ,omitempty" bson:"RJ,omitempty"`
	SK State  `json:"SK,omitempty" bson:"SK,omitempty"`
	TG State  `json:"TG,omitempty" bson:"TG,omitempty"`
	TN State  `json:"TN,omitempty" bson:"TN,omitempty"`
	TR State  `json:"TR,omitempty" bson:"TR,omitempty"`
	TT State  `json:"TT,omitempty" bson:"TT,omitempty"`
	UP State  `json:"UP,omitempty" bson:"UP,omitempty"`
	WB State  `json:"WB,omitempty" bson:"WB,omitempty"`
}

func GetIndiaEndpoint(response http.ResponseWriter, request *http.Request) {}

func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) {

	collection := client.Database("quickstart").Collection("covid")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	var err error
	var india bson.M
	if err = collection.FindOne(ctx, bson.M{}).Decode(&india); err != nil {
		log.Fatal(err)
	}
	// fmt.Println(podcast)

	response.Header().Set("content-type", "application/json")
	json.NewEncoder(response).Encode(india)

}

func Create(india India) {
	// response.Header().Set("content-type", "application/json")

	collection := client.Database("quickstart").Collection("covid")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, india)
	fmt.Println(result)
	// json.NewEncoder(response).Encode(result)
}

type Response struct {
	ID     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
}

func main() {

	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	database := client.Database("quickstart")
	covidCollection := database.Collection("covid")

	covidCollection.Drop(ctx)

	fmt.Println("Calling API...")
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://data.covid19india.org/v4/min/data.min.json", nil)
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	var responseObject India
	json.Unmarshal(bodyBytes, &responseObject)
	fmt.Printf("API Response as struct %+v\n", responseObject)

	Create(responseObject)

	router := mux.NewRouter()

	router.HandleFunc("/india", GetPeopleEndpoint).Methods("GET")
	http.ListenAndServe(":12345", router)
}
