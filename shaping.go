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

//Specification: just like Object-Pool mechanism, the first thing is to access a
// connection with the function : AccessConn(), when use over ,you should release
// the connection

package tshaping

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/kricen/tshaping/db"

	"github.com/kricen/tshaping/pool"
)

const (
	redisHashName = "tshaping_hash"
)

var (
	//ErrTimeout ： global error
	ErrTimeout = errors.New("timeout")
)

// TrafficShaping :
type TrafficShaping struct {
	burst       int              //default burst size, when beyond the size,wait until burst is not overflow
	shapingMap  map[string]int   // interface: max-burst
	connections *pool.ObjectPool // object pool
	mu          sync.RWMutex     // redis addon cann't performance well in concurrent circumstance, so need mutex to handle that
}

// idea : don't need to use redis ,can use sync.map[chan]interface{} to complete the
// limit funciton

// InitShaping : when use traffic shaping Algorithm init first
// because the achievement based on redis, need to registe a
// redis pool
func InitShaping(host, port, password string, maxBurst int, specialShaping map[string]int) (ts *TrafficShaping, err error) {
	// db.InitRedis()
	if maxBurst <= 0 {
		maxBurst = 600
	}
	if specialShaping == nil {
		specialShaping = make(map[string]int, 0)
	}
	ts = &TrafficShaping{burst: maxBurst, shapingMap: specialShaping}
	db.InitRedis(host, port, password)

	// init shaping environment : clear pre work-environment
	conn := db.GetRedisConn()
	_, err = conn.Do("del", redisHashName)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// init connection pool
	op := pool.NewPool()
	ts.connections = op
	return
}

// AccessConn : borrow a connection to user
func (t *TrafficShaping) AccessConn(url string) (pool *pool.Connection, err error) {
	err = t.accessConnFromRedis(url)
	if err != nil {
		return
	}
	pool = t.connections.AccessConnection(url)
	return
}

// CloseConn : return a connection which borrow from pool
func (t *TrafficShaping) CloseConn(conn *pool.Connection) (err error) {

	t.mu.Lock()
	defer t.mu.Unlock()

	redisConn := db.GetRedisConn()
	defer redisConn.Close()

	_, err = redisConn.Do("HINCRBY", redisHashName, conn.Name, -1)
	if err != nil {
		return
	}
	// release conn
	t.connections.ReleaseConnection(conn)

	return
}

// CloseShaping : release the resources when use over
func (t *TrafficShaping) CloseShaping() {
	// release the redis resource
	db.CloseRedis()
}

// AddUpdateShaping : add special shaping rule
func (t *TrafficShaping) AddUpdateShaping(specialShaping map[string]int) {
	for k, v := range specialShaping {
		t.shapingMap[k] = v
	}
}

func (t *TrafficShaping) accessConnFromRedis(url string) (err error) {
	var (
		onlineNumber int
	)
	t.mu.Lock()
	defer t.mu.Unlock()

	conn := db.GetRedisConn()
	defer conn.Close()

	for i := 0; i < 5; i++ {
		onlineNumber, err = redis.Int(conn.Do("HINCRBY", redisHashName, url, 1))
		if err != nil {
			fmt.Printf("Oops1:%s\n", err.Error())
			return
		}
		if _, ok := t.shapingMap[url]; !ok {
			t.shapingMap[url] = t.burst
		}

		if onlineNumber > t.shapingMap[url] {
			fmt.Printf("Oops11:%d,%d\n", onlineNumber, t.shapingMap[url])
			conn.Do("HINCRBY", redisHashName, url, -1)
			time.Sleep(10 * time.Millisecond)
			continue
		}
		return nil
	}

	return ErrTimeout
}
