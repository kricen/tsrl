package tshaping

import (
	"fmt"
	"os"
	"sync"
	"testing"
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
