package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	KDB "tigre/kdb"

	"context"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

var (
	keyDB    *redis.Client
	ctx      = context.Background()
	lastList KDB_List // redundant list copy
)

const (
	fruit_key = "fruits"
	gods_key  = "gods"
)

type TestData struct {
	Val string `json:"val"` // needs to be in caps to be exported
}

type KDB_List struct {
	Inputs TestData
	List   []string
}

func helloResp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome 🐯"))
}

// get the fruits list
func jsonFruits(w http.ResponseWriter, r *http.Request) {
	input := KDB.List2JSON(keyDB, fruit_key)
	w.Write(input)
}

func delGod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//origin := req.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == "DELETE" {

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
		_, n_times := KDB.List_DelStr(data.Val, "gods", keyDB)
		outputStr := fmt.Sprint(n_times)
		w.Write([]byte(outputStr))

	}
}

func postGod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
			err := json.Unmarshal(body, &data)
			if err != nil {
				fmt.Println(err)
			}
			result = KDB.List_Add(keyDB, "gods", data.Val)
		}
		outputStr := fmt.Sprintf("%s %s", data.Val, result)
		fmt.Println(outputStr)
		w.Write([]byte(result))
	}
}

func testGetList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	//params := mux.Vars(r)

	u, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing URL: %v", err), http.StatusBadRequest)
		return
	}

	// 2. Get the Query parameters
	query := u.Query()

	// 3. Access individual parameters
	val := query.Get("val")
	fmt.Println(val)

	w.Write([]byte(val))

	// pp, _ := url.ParseQuery(r.URL.RawQuery)

	// var key string
	// var val string
	// var data TestData

	// if r.Body != nil {
	// 	//fmt.Println("body:", r.Body)
	// 	ss := json.NewDecoder(r.Body).Decode(&data)
	// 	fmt.Println("Testdata:", data.Val, ss)
	// }

	// for k := range pp {
	// 	key = k
	// 	val = pp[k][0]
	// }
	// fmt.Println(key, val)
	// w.Write([]byte(key + ":" + val))
}

// this is the correct way to do it with the URL and GET request (not the body)
func jsonWildtype2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	//array holders
	var strArr []string
	var byteArr []byte
	var listName string

	// REFACTOR the KDB_list struct

	// test data holds the parsed JSON body for {"val":"listname"}
	var data TestData
	// KDB_list is nested
	var listObj KDB_List

	u, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing URL: %v", err), http.StatusBadRequest)
		return
	}

	// 2. Get the Query parameters
	query := u.Query()

	// 3. Access individual parameters
	val := query.Get("list")

	if val != "" {

		data.Val = val // TypeData
		listObj.Inputs = data
		listName = listObj.Inputs.Val

		strArr, byteArr = KDB.List2JSON_alpha(keyDB, listName)
		listObj.List = strArr
	}
	//save the last list requested
	if len(strArr) == 0 {
		outputStr := fmt.Sprint("❌ List not found: ", listName)
		w.Write([]byte(outputStr))
	} else {
		lastList = listObj
		w.Write(byteArr)
	}

}

func jsonWildtype(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	defer r.Body.Close()
	//array holders
	var strArr []string
	var byteArr []byte
	var listName string

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read body", http.StatusBadRequest)
		return
	}
	// test data holds the parsed JSON body for {"val":"listname"}
	var data TestData
	// KDB_list is nested
	var listObj KDB_List
	fmt.Println("request body:", string(body))

	if body != nil {
		err := json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println(err)
		}
		listObj.Inputs = data
		listName = listObj.Inputs.Val
		strArr, byteArr = KDB.List2JSON_alpha(keyDB, listObj.Inputs.Val)
		listObj.List = strArr
	}
	//save the last list requested
	if len(strArr) == 0 {
		outputStr := fmt.Sprint("❌ List not found: ", listName)
		w.Write([]byte(outputStr))
	} else {
		lastList = listObj
		w.Write(byteArr)
	}

}

// get the gods list
func jsonGods(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//origin := req.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	input := KDB.List2JSON(keyDB, gods_key)
	w.Write(input)
}

// simple post request --> /gods {"name":"Pazuzu"}
// LPUSH gods "Pazuzu"
// del

func main() {
	fmt.Println("hellos")

	keyDB = KDB.ConnectDB()

	// err := keyDB.LPush(ctx, "gods", "shebang").Err()
	// fmt.Println(err)

	//fruit_salad(keyDB, 10)

	r := mux.NewRouter()
	r.HandleFunc("/", helloResp).Methods("GET")
	r.HandleFunc("/t", testGetList).Methods("GET")
	r.HandleFunc("/jFruit", jsonFruits).Methods("GET")
	r.HandleFunc("/jGods", jsonGods).Methods("GET")
	r.HandleFunc("/anylist", jsonWildtype2).Methods("GET", "OPTIONS")
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
