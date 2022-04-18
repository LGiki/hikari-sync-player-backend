package gredis

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"hikari_sync_player/pkg/setting"
	"time"
)

var RedisConn *redis.Pool

func Setup() error {
	RedisConn = &redis.Pool{
		MaxIdle:     setting.GlobalSettings.Redis.MaxIdle,
		MaxActive:   setting.GlobalSettings.Redis.MaxActive,
		IdleTimeout: time.Duration(setting.GlobalSettings.Redis.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", setting.GlobalSettings.Redis.Host)
			if err != nil {
				return nil, err
			}
			if setting.GlobalSettings.Redis.Password != "" {
				if _, err := c.Do("AUTH", setting.GlobalSettings.Redis.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return nil
}

func Set(key string, data interface{}, time int) error {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}

	return nil
}

func Exists(key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

func Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

func LikeDeletes(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}

func PingRedis() error {
	conn := RedisConn.Get()
	defer conn.Close()
	_, err := conn.Do("PING")
	return err
}
