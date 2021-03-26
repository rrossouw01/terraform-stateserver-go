package main
// https://github.com/MerlinDMC/go-terraform-stateserver/blob/master/main.go
// v0.9.3
import (
    "flag"
    "io/ioutil"
    "strings"
    "log"
    "net/http"
    "context"
    "fmt"
    "encoding/json"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var (
	flagCertFile string
	flagKeyFile  string
	flagListenAddress string
)

func CheckError(e error) {
    if e != nil {
        fmt.Println(e)
    }
}

func init() {
	flag.StringVar(&flagCertFile, "certfile", "", "path to the certs file for https")
	flag.StringVar(&flagKeyFile, "keyfile", "", "path to the keyfile for https")
	flag.StringVar(&flagListenAddress, "listen_address", "0.0.0.0:8080", "address:port to bind listener on")
}

func main() {
	if !flag.Parsed() {
		flag.Parse()
	}

	router := http.NewServeMux()
	router.HandleFunc("/", requestHandler)

	if flagCertFile != "" && flagKeyFile != "" {
		http.ListenAndServeTLS(flagListenAddress, flagCertFile, flagKeyFile, router)
	} else {
		http.ListenAndServe(flagListenAddress, router)
	}
}

func requestHandler(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		http.NotFound(res, req)
		return
	}

	stateID := strings.Replace(req.URL.Path, "/", "",-1)
        clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
        client, e := mongo.Connect(context.TODO(), clientOptions)
        CheckError(e)

        e = client.Ping(context.TODO(), nil)
        CheckError(e)

        collection := client.Database("tfstatedb").Collection("states")

	switch req.Method {
	case "GET":
		fmt.Printf("received GET for state file: %s > ", stateID)

                var result interface{}

                docID, e := primitive.ObjectIDFromHex(stateID)
                CheckError(e)

                filter := bson.M{"_id": docID}
                projection := options.FindOne().SetProjection(bson.M{"_id": 0})
                tempResult := bson.M{}
                e = collection.FindOne(context.TODO(), filter, projection ).Decode(&tempResult)
                if e == nil {
                  obj, _ := json.Marshal(tempResult)
                  e = json.Unmarshal(obj, &result)
                } else {
		   goto not_found
                }
		fmt.Printf("FindOne = %s\n", stateID)
		res.WriteHeader(200)

                js, err := json.Marshal(result)
                if err != nil {
                  http.Error(res, err.Error(), http.StatusInternalServerError)
                  return
                }

                res.Header().Set("Content-Type", "application/json")
                res.Write(js)

		return
	case "POST":
		fmt.Printf("received POST for state file %s > ", stateID)
                // delete first
                pID, err := primitive.ObjectIDFromHex(stateID)
                result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": pID})
                if err != nil {
                  log.Fatal(err)
                }
                fmt.Printf("DeleteOne = %v > ", result)

		body, e := ioutil.ReadAll(req.Body)
		CheckError(e)
                var tf map[string]interface{}
                json.Unmarshal(body, &tf)
                oid, err := primitive.ObjectIDFromHex(stateID)
                tf["_id"] = oid

                _, e = collection.InsertOne(context.TODO(), tf)
                CheckError(e)
		fmt.Println("InsertOne = %v", result)
		
		return
	case "DELETE":
		fmt.Printf("received DELETE for state file: %s > ", stateID)
                pID, err := primitive.ObjectIDFromHex(stateID)
                result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": pID})
                if err != nil {
                  log.Fatal(err)
                }
		fmt.Printf("DeleteOne = %v\n", result)
        
		return
	}

not_found:
	http.NotFound(res, req)
}
