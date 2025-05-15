package KDB

// gatekeeper between Go & KeyDB server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	KeyDB *redis.Client
	ctx   = context.Background()
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

// LIST FUNCTIONS
//  _       _________ _______ _________ _______
// ( \      \__   __/(  ____ \\__   __/(  ____ \
// | (         ) (   | (    \/   ) (   | (    \/
// | |         | |   | (_____    | |   | (_____
// | |         | |   (_____  )   | |   (_____  )
// | |         | |         ) |   | |         ) |
// | (____/\___) (___/\____) |   | |   /\____) |
// (_______/\_______/\_______)   )_(   \_______)

// convert a textfile into a redis list
func Txt2List(rdb *redis.Client, keyname string, targetFile string) {
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
func List_RandItem(rdb *redis.Client, keyname string) string {
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

// add to a list (Lpush)
func List_Add(rdb *redis.Client, keyname string, val string) string { // 1 is success
	err := rdb.LPush(ctx, keyname, val).Err()
	if err != nil {
		fmt.Println("something failed", err)
		return "0"
	}
	fmt.Println(`"`, val, `"`, "-added to list-")
	return "1"
}

// del from a list (LRem)
func List_DelStr(key string, targetList string, kdb *redis.Client) (string, int) {
	result, err := kdb.Do(ctx, "LREM", targetList, 0, key).Result()
	i_result := int(result.(int64))
	fmt.Println("n removed:", i_result)
	if err != nil {
		fmt.Println("err:", err)
		return key, 0
	}
	return key, i_result

}

// return the json byte string
func List2JSON(rdb *redis.Client, listKey string) []byte {
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

// returns the string array too
func List2JSON_alpha(rdb *redis.Client, listKey string) ([]string, []byte) {
	var output []byte
	var strArr []string
	values, err := rdb.LRange(ctx, listKey, 0, -1).Result()
	if err != nil {
		fmt.Printf("Failed to read Redis list: %v \n", err)
		return strArr, output
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(values, "", "  ")
	if err != nil {
		fmt.Printf("Failed to parse JSON: %v \n", err)
		return strArr, output
	}
	return values, jsonData
}
