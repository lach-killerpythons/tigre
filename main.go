package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"

	"context"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

var (
	keyDB *redis.Client
	ctx   = context.Background()
)

const (
	fruit_key = "fruits"
	gods_key  = "gods"
)

func ConnectDB() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // host:port of the redis server
		Password: "",               // no password set
		DB:       0,                // use default DB
	})
	//panic if no db connected
	_, err := rdb.Ping(ctx).Result()
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
		rdb.LPush(ctx, keyname, line)

		fmt.Println("added to '", keyname, "' ~", line)
	}
}

// get a random list item
func KDB_rand_ListItem(rdb *redis.Client, keyname string) string {
	listLen, err := rdb.LLen(ctx, keyname).Result() //
	if err != nil {
		return err.Error()
	}
	random_i := rand.Intn(int(listLen)) // generate random index in range
	result, err := rdb.Do(ctx, "LRANGE", keyname, random_i, random_i).Result()
	if err != nil {
		return err.Error()
	}
	getItem := result.([]interface{})[0].(string)
	return getItem
}

func KDB_add2list(rdb *redis.Client, keyname string, val string) string { // 1 is success
	//result, err := rdb.Do("lpush", keyname, val).Result()

	//test_db := ConnectDB()

	//result, err := test_db.LPush(keyname, val).Result()
	err := keyDB.LPush(ctx, keyname, val).Err()
	if err != nil {
		fmt.Println("something failed", err)
		return "0"
	}
	fmt.Println(`"`, val, `"`, "-added to list-")
	return "1"
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

// get the fruits list
func jsonFruits(w http.ResponseWriter, r *http.Request) {
	input := KDB_list2JSON(keyDB, fruit_key)
	w.Write(input)
}

type TestData struct {
	Val string `json:"val"` // needs to be in caps to be exported
}

func KDB_del_str(key string, targetList string, kdb *redis.Client) (string, int) {
	//kdb.LRem(ctx, targetList, 0, key)
	result, err := kdb.Do(ctx, "LREM", targetList, 0, key).Result()
	i_result := int(result.(int64))
	fmt.Println("n removed:", i_result)
	if err != nil {
		fmt.Println("err:", err)
		return key, 0
	}
	return key, i_result

}

func delGod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read body", http.StatusBadRequest)
		return
	}
	var data TestData

	if body != nil {
		//err := json.NewDecoder(body).Decode(&data)
		err := json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println(err)
		}
	}
	key, n_times := KDB_del_str(data.Val, "gods", keyDB)
	outputStr := fmt.Sprintf(key, "N removed:", n_times)
	w.Write([]byte(outputStr))
}

func postGod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//origin := req.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	if r.Method == "POST" {
		defer r.Body.Close()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Cannot read body", http.StatusBadRequest)
			return
		}
		var result string
		var data TestData
		fmt.Println("request body:", string(body))

		if body != nil {
			//err := json.NewDecoder(body).Decode(&data)
			err := json.Unmarshal(body, &data)
			if err != nil {
				fmt.Println(err)
			}
			result = KDB_add2list(keyDB, "gods", data.Val)
		}
		outputStr := fmt.Sprintf("%s %s", data.Val, result)
		fmt.Println(outputStr)
		w.Write([]byte(result))
	}
}

func testPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	//params := mux.Vars(r)
	pp, _ := url.ParseQuery(r.URL.RawQuery)

	var key string
	var val string
	var data TestData

	if r.Body != nil {
		//fmt.Println("body:", r.Body)
		ss := json.NewDecoder(r.Body).Decode(&data)
		fmt.Println("Testdata:", data.Val, ss)
	}

	for k := range pp {
		key = k
		val = pp[k][0]
	}
	fmt.Println(key, val)
	w.Write([]byte(key + ":" + val))
}

// get the gods list
func jsonGods(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//origin := req.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	input := KDB_list2JSON(keyDB, gods_key)
	w.Write(input)
}

func KDB_list2JSON(rdb *redis.Client, listKey string) []byte {
	var output []byte
	values, err := rdb.LRange(ctx, listKey, 0, -1).Result()
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

// simple post request --> /gods {"name":"Pazuzu"}
// LPUSH gods "Pazuzu"
// del

func main() {
	fmt.Println("hellos")

	keyDB = ConnectDB()

	// err := keyDB.LPush(ctx, "gods", "shebang").Err()
	// fmt.Println(err)

	//fruit_salad(keyDB, 10)

	r := mux.NewRouter()
	r.HandleFunc("/", helloResp).Methods("GET")
	r.HandleFunc("/jFruit", jsonFruits).Methods("GET")
	r.HandleFunc("/jGods", jsonGods).Methods("GET")
	r.HandleFunc("/new_god", postGod).Methods("POST", "OPTIONS")
	r.HandleFunc("/del", delGod).Methods("DELETE", "OPTIONS")
	srv := &http.Server{
		Addr:    "localhost:8080",
		Handler: r,
	}
	srv.ListenAndServe()

	// 	36) "Zeus"
	// 37) "Athena"
	// 38) "Apollo"
	// 39) "Artemis"
	// 40) "Dionysus"
	// 41) "Serapis"
	// 42) "Hermes"
	// 43) "Demeter"
	// 44) "Aphrodite"
	// 45) "Asclepius"

}
