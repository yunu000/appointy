package main
//necessary imports
import (
	// "strconv" //for converting string into int
	"time" //for timestamp
	"fmt" //for console output
	"log" //logging errors
	"net/http" //for routing
	"github.com/gorilla/mux" //for handing routing variable/wildcard
	"context"
	"go.mongodb.org/mongo-driver/mongo" //for connecting mongodb
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
 )

//creating a global variable of type *mongo.Client
var client *mongo.Client

//function to create a global copy of client
func passClient(c *mongo.Client) {
	client=c
}
func getUser(id string,c string) []bson.M {

	collection := client.Database("appointy").Collection(c)

	cursor,err := collection.Find(context.TODO(),bson.M{"id":id})

	if err != nil {
		log.Fatal(err)
	}
	var users []bson.M
	if err := cursor.All(context.TODO(),&users);err!= nil{
		log.Fatal(err)
	}
	return users
}
// Function for post route /user
func userRoute(w http.ResponseWriter, r *http.Request){
	if err:=r.ParseForm(); err!=nil{
		fmt.Fprintf(w,"ParseForm() error : %v",err)
	}
	//Get data from request
	id := r.FormValue("id")
	name := r.FormValue("name")
	email := r.FormValue("email")
	pno:= r.FormValue("pno")
	dob :=r.FormValue("dob")
	// Create instance of collection
	collection := client.Database("appointy").Collection("users")
	// Insert into collection in bson format
	insertResult, err := collection.InsertOne(context.TODO(),bson.D{
		{"id",id},
		{"name",name},
		{"dob",dob},
		{"pno",pno},
		{"email",email},
	})
	if err!=nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted post with ID:", insertResult.InsertedID)
	fmt.Fprintf(w,"Successfuly Inserted")
}
// Function for get route /users/user_id
func usersRoute(w http.ResponseWriter, r *http.Request){
	vars:=mux.Vars(r)
	varid:=vars["id"]
	user:=getUser(varid,"users")
	fmt.Fprintf(w,"%v",user)
}
//Function for route /contact and /contact?user=<user id>&infection_timestamp=<timestamp>’
func contactRoute(w http.ResponseWriter, r *http.Request){
	if r.Method=="GET"{
		keys1, _ := r.URL.Query()["user"]
		keys2, _ := r.URL.Query()["infection_timestamp"]
		key1:=keys1[0]
		key2:=keys2[0]
		timestamp, _:=time.Parse("2006-01-02 3:4:5",key2)
		fmt.Println(timestamp)
		contacts:=getUser(key1,"contact")
		for _,contact := range contacts{
			contact_date, _:=time.Parse("2006-01-02 3:4:5",fmt.Sprint(contact["timestamp"]))
			contact_date.Add(time.Hour * 24 * 7 * time.Duration(2))
			fmt.Println(contact_date)
		}
		fmt.Fprintf(w,"%v",contacts)
		return 
	}
	if err:=r.ParseForm(); err!=nil{
		fmt.Fprintf(w,"ParseForm() error : %v",err)
	}
	id1 := r.FormValue("id1")
	id2 :=r.FormValue("id2")
	currentTime :=time.Now().Format("2006-01-02 3:4:5")
	collection := client.Database("appointy").Collection("contact")
	insertResult, err := collection.InsertOne(context.TODO(),bson.D{
		{"id",id1},
		{"id1",id2},
		{"timestamp", currentTime},
	})
	if err!=nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted post with ID:", insertResult.InsertedID)
	fmt.Fprintf(w,"Successfuly posted")
}
//main function
func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
	    log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
	    log.Fatal(err)
	}
	passClient(client)
	fmt.Println("Connected to MongoDB+!")
	// mongodb connection
	// GetPost()
	rtr := mux.NewRouter()
	http.HandleFunc("/index",indexRoute)
	http.HandleFunc("/user",userRoute)
	http.HandleFunc("/contacts",contactRoute)
	rtr.HandleFunc("/users/{id:[0-9]+}",usersRoute)
	http.Handle("/", rtr)
	err1:=http.ListenAndServe(":8000",nil)
	if err1 !=nil{
		log.Fatal(err1)
	}
}