package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

var (
	keyDB *redis.Client
)

const (
	fruit_key = "fruits"
)

func ConnectDB() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // host:port of the redis server
		Password: "",               // no password set
		DB:       0,                // use default DB
	})
	//panic if no db connected
	_, err := rdb.Ping().Result()
	if err != nil {
		panic(err)
	}
	return rdb
}

// convert a textfile into a redis list
func KDB_File2List(rdb *redis.Client, keyname string) {
	targetFile := "fruity.txt"
	file, err := os.Open(targetFile)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		rdb.LPush(keyname, line)

		fmt.Println("added to '", keyname, "' ~", line)
	}
}

// get a random list item
func KDB_rand_ListItem(rdb *redis.Client, keyname string) string {
	listLen, err := rdb.LLen(keyname).Result() //
	if err != nil {
		return err.Error()
	}
	random_i := rand.Intn(int(listLen)) // generate random index in range
	result, err := rdb.Do("LRANGE", keyname, random_i, random_i).Result()
	if err != nil {
		return err.Error()
	}
	getItem := result.([]interface{})[0].(string)
	return getItem
}

func fruit_salad(rdb *redis.Client, n_fruits int) []string {
	var fruits []string
	for i := 0; i < n_fruits; i++ {
		next_fruit := KDB_rand_ListItem(rdb, "fruits")
		fruits = append(fruits, next_fruit)
	}
	ss := ""
	for _, s := range fruits {
		ss += s + ", "
	}
	fmt.Println("today i made a salad containing: ", ss)
	return fruits
}

func helloResp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome 🐯"))
}

func jsonFruits(w http.ResponseWriter, r *http.Request) {
	input := getFruityJSON(keyDB, fruit_key)
	w.Write(input)
}

func getFruityJSON(rdb *redis.Client, listKey string) []byte {
	var output []byte
	values, err := rdb.LRange(listKey, 0, -1).Result()
	if err != nil {
		fmt.Printf("Failed to read Redis list: %v \n", err)
		return output
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(values, "", "  ")
	if err != nil {
		fmt.Printf("Failed to parse JSON: %v \n", err)
		return output
	}
	return jsonData
}

func main() {
	fmt.Println("hellos")
	keyDB = ConnectDB()
	fruit_salad(keyDB, 10)

	r := mux.NewRouter()
	r.HandleFunc("/", helloResp).Methods("GET")
	r.HandleFunc("/jFruit", jsonFruits).Methods("GET")
	srv := &http.Server{
		Addr:    "localhost:8080",
		Handler: r,
	}
	srv.ListenAndServe()
}
