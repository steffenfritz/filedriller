package filedriller

import (
	"github.com/gomodule/redigo/redis"
)

// RedisConnect creates a connection to a Redis server
func RedisConnect(r RedisConf) redis.Conn {
	conn, err := redis.Dial("tcp", *r.Server+":"+*r.Port)
	if err != nil {
		e(err)
	}

	return conn
}

// RedisGet returns the boolean answering if a hash is in the NSRL
func RedisGet(conn redis.Conn, hashSum string) string {
	inNSRL, err := redis.String(conn.Do("GET", hashSum))
	if err == redis.ErrNil {
		inNSRL = "FALSE"
	}

	return inNSRL
}
