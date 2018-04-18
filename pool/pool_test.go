package pool

import (
	"testing"
	"time"
)

var op *ObjectPool

func init() {
	op = NewPool()
}

func TestAccessConnection(t *testing.T) {
	timeout := time.After(100 * time.Second)
	tick := time.Tick(1 * time.Millisecond)

	for {
		select {
		case <-timeout:
			return
		case <-tick:
			conn := op.AccessConnection("hello")
			go func() {
				time.Sleep(3 * time.Millisecond)
				op.ReleaseConnection(conn)
			}()
		}

	}
}

func TestAccessAndReleaseConnetion(t *testing.T) {
	for i := 0; i < 2000000000; i++ {
		conn := op.AccessConnection("hello")
		go func() {
			time.Sleep(1 * time.Second)
			//fmt.Println(&conn)
			op.ReleaseConnection(conn)
		}()
	}
	//	time.Sleep(30 * time.Second)
}
