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
	"github.com/kricen/tshaping/db"
	"github.com/kricen/tshaping/pool"
)

const (
	redisHashName = "tshaping_hash"
)

// TrafficShaping :
type TrafficShaping struct {
	burst       int            //default burst size, when beyond the size,wait until burst is not overflow
	shapingMap  map[string]int // interface: max-burst
	connections *pool.ObjectPool
}

// InitShaping : when use traffic shaping Algorithm init first
// because the achievement based on redis, need to registe a
// redis pool
func InitShaping(serviceName, host, port, password string, maxBurst int, specialShaping map[string]int) (ts *TrafficShaping, err error) {
	// db.InitRedis()
	ts = &TrafficShaping{burst: maxBurst, shapingMap: specialShaping}
	db.InitRedis(serviceName, host, port, password)

	// init shaping environment : clear pre work-environment
	conn := db.GetRedisConn()
	conn.Do("del", redisHashName)
	defer conn.Close()

	// init connection pool
	op := pool.NewPool()
	ts.connections = op
	return
}

// AccessConn : borrow a connection to user
func (t *TrafficShaping) AccessConn(url string) *pool.Connection {
	// Notice : maybe occur a codition : cache-breakdown

	return nil
}

// CloseConn : return a connection which borrow from pool
func (t *TrafficShaping) CloseConn(conn *pool.Connection) {

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
