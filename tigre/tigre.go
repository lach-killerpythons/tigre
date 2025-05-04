package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"

	"github.com/go-redis/redis"
)

var ctx = context.Background()

func kdb_File2List(rdb *redis.Client, keyname string) {
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

func kdb_rand_ListItem(rdb *redis.Client, keyname string) string {
	listLen, err := rdb.LLen(keyname).Result() //
	if err != nil {
		return err.Error()
	}
	random_i := rand.Intn(int(listLen)) // generate random index in range
	//cmdStr := fmt.Sprintf(`LRANGE fruits %d, %d`, random_i, random_i)
	//fmt.Println(cmdStr)
	result, err := rdb.Do("LRANGE", keyname, random_i, random_i).Result()
	if err != nil {
		return err.Error()
	}

	// extract interface index and force to string
	getItem := result.([]interface{})[0].(string)

	//return reflect.ValueOf(getItem).String()
	//return fmt.Sprintf(string(getItem))
	return getItem
}

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

	// ping
	// pong, err := rdb.Ping().Result()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(pong)
	// // Randomkey
	// randomkey, _ := rdb.RandomKey().Result()
	// fmt.Println(randomkey)
	// // KEYS *
	// keyslice, _ := rdb.Do("KEYS", "*").Result()
	// fmt.Println(keyslice)
	// // Get
	// val, err := rdb.Get("tigre").Result()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Welcome to", val)

	// // lpush (Listadd)
	// rdb.LPush("list0", "mandarin").Result()

	// // lrange Getlist
	// getList, _ := rdb.LRange("list0", 0, 10).Result()
	// fmt.Println(getList)
}

func fruit_salad(rdb *redis.Client, n_fruits int) []string {
	var fruits []string
	for i := 0; i < n_fruits; i++ {
		next_fruit := kdb_rand_ListItem(rdb, "fruits")
		fruits = append(fruits, next_fruit)
	}
	ss := ""
	for _, s := range fruits {
		ss += s + ", "
	}
	fmt.Println("today i made a salad containing: ", ss)
	return fruits
}

func main() {
	//connectDB()

	keyDB := ConnectDB()
	//kdb_File2List(keyDB, "fruits")
	fruit_salad(keyDB, 10)
}
