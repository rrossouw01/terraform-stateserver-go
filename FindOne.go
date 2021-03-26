package main
 
import (
    "flag"
    "log"
    "context"
    "fmt"
    "encoding/json"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var (
     flagId string
)
 
func init() {
     flag.StringVar(&flagId, "id", "", "mongodb ObjectId")
}

func PrettyPrint(v interface{}) (err error) {
      b, err := json.MarshalIndent(v, "", "  ")
      if err == nil {
              fmt.Println(string(b))
      }
      return
}
 
func main() {
    __version := "v0.9.2"
    fmt.Println("FindOne.go " + __version)
    // run example: ~/tf-server-mongodb-go$ go run FindOne.go -id=605b593e5e5cce2317e4a552
    //605a7c438f8e57850529ff9c, 605b593e5e5cce2317e4a552, 605b5f6de262e8d9665baaed
    
    if !flag.Parsed() {
	flag.Parse()
    }
	
    if flagId != "" {
          fmt.Println("finding id: " + flagId)
	} else {
          log.Fatal("no _id specified")
	}

    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, e := mongo.Connect(context.TODO(), clientOptions)
    CheckError(e)
 
    e = client.Ping(context.TODO(), nil)
    CheckError(e)
 
    collection := client.Database("tfstatedb").Collection("states")
 
    var res interface{}

    docID, e := primitive.ObjectIDFromHex(flagId)
    CheckError(e)
    
    filter := bson.M{"_id": docID}
    projection := options.FindOne().SetProjection(bson.M{"_id": 0})
    tempResult := bson.M{}
    //err1 := collection.FindOne(context.TODO(), filter).Decode(&tempResult)
    //e = collection.FindOne(context.TODO(), filter ).Decode(&tempResult)
    e = collection.FindOne(context.TODO(), filter, projection ).Decode(&tempResult)
    if e == nil {   
      obj, _ := json.Marshal(tempResult)
      e = json.Unmarshal(obj, &res)
    }
    PrettyPrint(res)
}
 
func CheckError(e error) {
    if e != nil {
        fmt.Println(e)
    }
}
