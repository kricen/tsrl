package tsrl

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	fmt.Println("start")
	pool := New()
	for i := 0; i < 100000; i++ {
		pool.GetBucket(fmt.Sprintf("tsrl%d", 1))
	}
	fmt.Println("end")
	fmt.Println("sleep 30 sec")
	fmt.Print(len(pool.bmap.Items()))
}
