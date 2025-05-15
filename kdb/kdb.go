package KDB

import (
	"bufio"
	"context"
	"fmt"
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
