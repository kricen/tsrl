package tsrl

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kricen/tsrl/model"
)

func TestNew(t *testing.T) {
	fmt.Println("start")
	pool := New()
	for i := 0; i < 100000; i++ {
		pool.GetBucket(fmt.Sprintf("tsrl%d", 1))
	}
	fmt.Println("end")
	fmt.Print(len(pool.bmap.Items()))
}

func TestFlow(t *testing.T) {
	var wg sync.WaitGroup
	var failedNum int64
	pool := New()
	// bucket := pool.AddBucket("hello", 1000, 5*time.Second, model.BUCKET_TYPE_RATE_LIMIT)
	bucket := pool.AddBucket("hello", 500, 5*time.Second, model.BUCKET_TYPE_TRAFFIC_SHAPING)

	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			token, err := bucket.BorrowToken(0)
			if err != nil {
				atomic.AddInt64(&failedNum, 1)
				return
			}
			time.Sleep(50 * time.Millisecond)
			bucket.ReleaseToken(token)
		}(i)
	}

	wg.Wait()
	fmt.Printf("failedNum:%d\n", failedNum)
}
