package tshaping

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

var (
	ts  *TrafficShaping
	err error
)

func init() {
	ts, err = InitShaping("", "", "", 20, nil)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

func TestShaping(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		for j := 0; j < 30; j++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				conn, err := ts.AccessConn(fmt.Sprintf("url%d", i))
				if err != nil {
					fmt.Printf("Oops:%s\n", err.Error())
					return
				}
				fmt.Printf("success\n")
				ts.CloseConn(conn)
			}(i)
		}
	}

	wg.Wait()
}

func TestChannel(t *testing.T) {
	var mu sync.RWMutex
	numm := make(map[int]int)
	tc := make(chan int)
	for i := 0; i < 10; i++ {
		go func(index int) {
			for {
				select {
				case v := <-tc:
					mu.Lock()
					numm[index]++
					mu.Unlock()
					fmt.Printf("from index:%d,value:%d \n", index, v)
				}
			}
		}(i)
	}

	for i := 0; i < 100000; i++ {
		tc <- i
		// time.Sleep(10 * time.Millisecond)
	}
	for k, v := range numm {
		fmt.Printf("consumer:%d,consume:%d\n", k, v)
	}

	time.Sleep(2 * time.Second)
}
