package pool

import (
	"fmt"
	"sync"
	"time"
)

//Connection : borrow a connection to user ,when finish release by user
type Connection struct {
	name       string
	updateTime int64
}

// ObjectPool : connection's pool
type ObjectPool struct {
	connections map[*Connection]bool
	rw          sync.RWMutex
}

// NewPool : generate a object pool
func NewPool() *ObjectPool {
	conns := make(map[*Connection]bool, 0)
	op := &ObjectPool{connections: conns}
	go op.recycling()
	return op
}

// AccessConnection : get a connection from pool
func (o *ObjectPool) AccessConnection(url string) *Connection {
	o.rw.RLock()
	defer o.rw.RUnlock()
	for k, v := range o.connections {
		if !v {
			o.connections[k] = true
			k.name = url
			k.updateTime = time.Now().Unix()
			return k
		}
	}
	conn := &Connection{name: url, updateTime: time.Now().Unix()}
	o.connections[conn] = true
	return conn
}

// ReleaseConnection : release connection
func (o *ObjectPool) ReleaseConnection(conn *Connection) {
	o.rw.RLock()
	defer o.rw.RUnlock()
	if _, ok := o.connections[conn]; ok {
		fmt.Println("存在啊啊啊啊")
		o.connections[conn] = false
	} else {
		fmt.Println("不存在啊啊啊 啊 ")
	}

}

// recycling : recycling the idle connection
func (o *ObjectPool) recycling() {
	tick := time.Tick(15 * time.Second)
	for {
		select {
		case <-tick:
			o.rw.Lock()
			defer o.rw.Unlock()
			timeNow := time.Now().Add(-120 * time.Second).Unix()
			for k, v := range o.connections {
				if !v {
					if k.updateTime > timeNow {
						delete(o.connections, k)
					}
				}
			}
		}

	}
}
