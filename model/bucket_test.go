package model

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func BenchmarkNew(t *testing.B) {
	bucket := New(100, 5000*time.Millisecond, BUCKET_TYPE_TRAFFIC_SHAPING)
	var wg sync.WaitGroup
	var failedNum int64
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			token, err := bucket.BorrowToken(0)
			if err != nil {
				atomic.AddInt64(&failedNum, 1)
				//fmt.Printf("Oops,work %d is Err: %s\n", i, err.Error())
				return
			}
			//time.Sleep(50 * time.Millisecond)
			//fmt.Println("fine,comelete the work:", i)
			bucket.ReleaseToken(token)
		}(i)
	}

	wg.Wait()
	// fmt.Printf("maxSize:%d,vacantSize:%d,undistributedSize:%d\n", bucket.maxSize, bucket.vacantSize, bucket.undistributedSize)
	fmt.Printf("failedNum:%d,consume:%d ,product :%d\n", failedNum, bucket.consume, bucket.produceToken)

}

func TestNewRLBucket(t *testing.T) {
	bucket := New(1000, 5*time.Second, BUCKET_TYPE_RATE_LIMIT)
	//time.Sleep(1 * time.Second)
	for i := 0; i < 10; i++ {
		go func() {
			bucket.BorrowToken(0)
			bucket.BorrowToken(0)
			bucket.BorrowToken(0)
			bucket.BorrowToken(0)
			bucket.BorrowToken(0)
			bucket.BorrowToken(0)
			bucket.BorrowToken(0)
			bucket.BorrowToken(0)
			bucket.BorrowToken(0)
			bucket.BorrowToken(0)
			bucket.BorrowToken(0)
			bucket.BorrowToken(0)

		}()
		fmt.Println(i)
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("consume:%d ,product :%d", bucket.consume, bucket.produceToken)
}
