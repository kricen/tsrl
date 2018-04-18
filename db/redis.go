// Copyright 2018 oliver kricen
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
//
//
//  go get github.com/gomodule/redigo/redis
//
//

package db

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

// register the redis pool
func newPool(server string, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
				}
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				c.Close()
			}
			return nil
		},
	}
}

// InitRedis : init redis pool with some param
func InitRedis(svcName, host, port, password string) {
	redishost := getValue(host, "127.0.0.1")
	redisport := getValue(port, "6379")
	redispwd := getValue(password, "")
	redis := fmt.Sprintf("%s:%s", redishost, redisport)
	setupRedis(redis, redispwd)
}

func setupRedis(server string, password string) {
	pool = newPool(server, password)
	err := pool.TestOnBorrow(pool.Get(), time.Now()) // time.Now() this parameter no use

	if err != nil {

	}
}

//CloseRedis : close the redis resources when complete the program
func CloseRedis() error {
	if pool != nil {
		return pool.Close()
	}
	return nil
}

// GetRedisConn : access a connectin from redis pool
func GetRedisConn() redis.Conn {
	return pool.Get()
}

func getValue(arg1, arg2 string) string {
	if arg1 == "" {
		return arg2
	}
	return arg1
}
