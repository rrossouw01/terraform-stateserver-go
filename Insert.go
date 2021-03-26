package main
 
import (
    "context"
    "fmt"
    "os"
    "io/ioutil"
    "encoding/json"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    //"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)
 
func main() {
    //run example: ~/tf-server-mongodb-go$ go run Insert.go
    
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, e := mongo.Connect(context.TODO(), clientOptions)
    CheckError(e)
 
    e = client.Ping(context.TODO(), nil)
    CheckError(e)
 
    collection := client.Database("tfstatedb").Collection("states")
 
    jsonFile, err := os.Open("terraform.tfstate")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("Successfully opened Terraform State File for Insert")
    defer jsonFile.Close()

    //fmt.Println("Inserting lineage: " + fmt.Sprintf("%v", result["lineage"]))

    byteValue, _ := ioutil.ReadAll(jsonFile)

    var tf map[string]interface{}
    json.Unmarshal([]byte(byteValue), &tf) 
    i := primitive.NewObjectID()   
    tf["_id"] = i

    _, e = collection.InsertOne(context.TODO(), tf)
    //_, e = collection.InsertOne(context.TODO(), {"_id": i, tf})
    CheckError(e)
    
    fmt.Println("InsertOne() _id:", i)
}
 
func CheckError(e error) {
    if e != nil {
        fmt.Println(e)
    }
}
