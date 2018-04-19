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

package pool

import (
	"sync"
	"time"
)

//Connection : borrow a connection to user ,when finish release by user
type Connection struct {
	Name       string
	updateTime int64
}

// ObjectPool : connection's pool
type ObjectPool struct {
	connections []*Connection
	rw          sync.RWMutex
	// gnum        int64
	// anum        int64
	// total       int64
}

// NewPool : generate a object pool
func NewPool() *ObjectPool {
	conns := make([]*Connection, 0)
	op := &ObjectPool{connections: conns}
	// go op.recycling()
	return op
}

// AccessConnection : get a connection from pool
func (o *ObjectPool) AccessConnection(url string) *Connection {
	o.rw.Lock()
	defer o.rw.Unlock()
	// atomic.AddInt64(&o.total, 1)
	if len(o.connections) > 0 {
		// atomic.AddInt64(&o.anum, 1)
		v := o.connections[0]
		// fmt.Println("reuse:", &v)
		v.Name = url
		v.updateTime = time.Now().Unix()
		o.connections = o.connections[1:]
		return v
	}
	// atomic.AddInt64(&o.gnum, 1)
	conn := &Connection{Name: url, updateTime: time.Now().Unix()}
	return conn
}

// ReleaseConnection : release connection
func (o *ObjectPool) ReleaseConnection(conn *Connection) {
	o.rw.Lock()
	defer o.rw.Unlock()
	// fmt.Println("release:", &conn)
	o.connections = append(o.connections, conn)
}

// recycling : recycling the idle connection
func (o *ObjectPool) recycling() {
	tick := time.Tick(1 * time.Second)
	for {
		select {
		case <-tick:
			// fmt.Println("reuse num:", o.anum, " produce num:", o.gnum, " total:", o.total)
		}
	}
}
