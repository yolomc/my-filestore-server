package redis

import (
	"log"

	"github.com/garyburd/redigo/redis"
)

//IsExists 判断key是否存在
func IsExists(key string) (bool, error) {
	redisConn := pool.Get()
	defer redisConn.Close()
	isExit, err := redis.Bool(redisConn.Do("EXISTS", key))
	if err != nil {
		log.Println("pkg.redis.IsExists Do error is:" + err.Error())

		return false, err
	}
	return isExit, nil
}

//Set Set
func Set(key, value string, ttl int) error {
	redisConn := pool.Get()
	defer redisConn.Close()

	_, err := redisConn.Do("SET", key, value, "EX", ttl)
	if err != nil {
		log.Println("pkg.redis.Set Do error is:" + err.Error())
		return err
	}
	return nil
}

//Get Get
func Get(key string) (string, error) {
	redisConn := pool.Get()
	defer redisConn.Close()

	value, err := redis.String(redisConn.Do("GET", key))
	if err != nil {
		log.Println("pkg.redis.Get Do error is:" + err.Error())
		return "", err
	}
	return value, nil
}

//Del Del
func Del(key string) error {
	redisConn := pool.Get()
	defer redisConn.Close()

	_, err := redisConn.Do("DEL", key)
	if err != nil {
		log.Println("pkg.gredis.Del Do error is:" + err.Error())
		return err
	}
	return nil
}
