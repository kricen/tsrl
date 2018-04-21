package bucket

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func BenchmarkNew(t *testing.B) {
	bucket := New(100, 5*time.Second)
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			token, err := bucket.BorrowToken(0)
			if err != nil {
				fmt.Printf("Oops,work %d is Err: %s\n", i, err.Error())
				return
			}
			time.Sleep(200 * time.Millisecond)
			fmt.Println("fine,comelete the work:", i)
			bucket.ReleaseToken(token)
		}(i)
	}

	wg.Wait()
	fmt.Printf("maxSize:%d,vacantSize:%d,undistributedSize:%d\n", bucket.maxSize, bucket.vacantSize, bucket.undistributedSize)

}
