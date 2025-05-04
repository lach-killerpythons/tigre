package api

import "github.com/go-redis/redis"

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
