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

//HSet （可批量）参数格式："key","field1","value1","field2","value2"...
func HSet(args ...interface{}) {
	redisConn := pool.Get()
	defer redisConn.Close()
	_, err := redisConn.Do("HMSET", args...)
	if err != nil {
		log.Println("pkg.gredis.HSet Do error is:", err)
		return
	}
	return
}

//HGet （可批量）参数格式："key","field1","field2"...
func HGet(key string, field interface{}) ([]interface{}, error) {
	redisConn := pool.Get()
	defer redisConn.Close()
	values, err := redis.Values(redisConn.Do("HMGET", key, field))
	if err != nil {
		log.Println("pkg.gredis.HGet Do error is:", err)
		return nil, err
	}
	return values, nil
}

//HGetAll 获取所有field和value
func HGetAll(key string) ([]interface{}, error) {
	redisConn := pool.Get()
	defer redisConn.Close()
	values, err := redis.Values(redisConn.Do("HGetAll", key))
	if err != nil {
		log.Println("pkg.gredis.HGetAll Do error is:", err)
		return nil, err
	}
	return values, nil
}

// func SetBit(key string, field, value interface{}) {
// 	redisConn := RedisCli.Get()
// 	defer redisConn.Close()
// 	_, err := redisConn.Do("SETBIT", key, field, value)
// 	if err != nil {
// 		log.Error("pkg.gredis.SETBIT Do error is:", err)
// 		return
// 	}
// 	return
// }

// func GetBit(key string, field interface{}) (interface{}, error) {
// 	redisConn := RedisCli.Get()
// 	defer redisConn.Close()
// 	reply, err := redisConn.Do("GetBit", key, field)
// 	if err != nil {
// 		log.Error("pkg.gredis.HGet Do error is:", err)
// 		return nil, err
// 	}
// 	return reply, nil
// }
