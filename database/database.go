package database

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

func CreateClient(dbNo int)  *redis.Client{
	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:  os.Getenv("DB_ADDRESS"),
	// 	Password: os.Getenv("PASSWORD"),
	// 	DB: dbNo,
	// })

	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(opt)

	_, _err := rdb.Ping(Ctx).Result()
    if _err != nil {
        log.Fatalln("Redis connection was refused")
    }

	return rdb
}