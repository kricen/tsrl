package db

import (
	"fmt"
	"sync"
	"testing"

	"github.com/garyburd/redigo/redis"
)

var (
	conn redis.Conn
	sm   sync.RWMutex
)

const (
	hashName = "mytestHash"
)

func initRedis() {
	InitRedis("", "", "")

}

func TestDel(t *testing.T) {
	initRedis()
	conn = GetRedisConn()
	defer conn.Close()
	v, err := redis.Int(conn.Do("del", hashName)) // not exist in redis
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log(v)
}

func BenchmarkAddUpdate(t *testing.B) {
	initRedis()
	var wg sync.WaitGroup
	// hash set or update
	fmt.Println()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			incr(hashName, "url1")
		}()

	}
	wg.Wait()
	fmt.Println()
}

func incr(hashName, url string) {
	sm.Lock()
	defer sm.Unlock()
	conn = GetRedisConn()
	v, err := redis.Int(conn.Do("HINCRBY", hashName, "url1", 1))
	if err != nil {
		conn.Close()
		fmt.Println("oooooo", err.Error())
		return
	}
	conn.Close()
	fmt.Print(v, ",")
}
