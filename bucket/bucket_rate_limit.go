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

package bucket

import (
	"fmt"
	"sync/atomic"
	"time"
)

// RLBucket : rate limit bucket
type RLBucket struct {
	qps               int64
	maxSize           int64
	vacantSize        int64
	undistributedSize int64
	RWChan            chan interface{}
	timeoutDuration   time.Duration
	Bucket
}

//NewRLBucket : init a traffic-shaping Bucket with two param : maxSize  and timeoutDuration
//      maxSize : max size that this bucket can hold
//      timeoutDuration : max time that borrow a token from bucket
func NewRLBucket(qps int64, timeoutDuration time.Duration) *RLBucket {
	if qps <= 0 {
		qps = defaultMaxSize // 10000
	}
	if timeoutDuration <= 0 {
		timeoutDuration = defaultTimeout
	}
	channel := make(chan interface{}, qps)
	bucket := &RLBucket{qps: qps, RWChan: channel, vacantSize: 0, timeoutDuration: timeoutDuration, maxSize: qps}
	go bucket.produce()
	return bucket
}

// BorrowToken : access token from bucket.
// Specification : if timeoutDruation is little than zero , program will set timeoutDruation equals Buckte's
//                  timeoutDruation set in init , while  if not set when function init, timeoutDruation is
//                  a default value : 3 seconds
func (b *RLBucket) BorrowToken(timeoutDruation time.Duration) (token interface{}, err error) {

	if timeoutDruation <= 0 {
		timeoutDruation = b.timeoutDuration
	}
	// define a received-only timeout channel
	timeout := time.After(timeoutDruation)

	select {
	case <-timeout:
		return nil, ErrTimeout
	case token = <-b.RWChan:
		fmt.Println("consume")
		atomic.AddInt64(&b.vacantSize, -1)
	}
	return
}

//ReleaseToken : don't need to exec
func (b *RLBucket) ReleaseToken(token interface{}) {
	// do nothing
}

// produce : a function to produce the token ,and throw token into bucket
func (b *RLBucket) produce() {

	// calcute the produce speed
	speed := int(float64(1) / float64(b.qps) * 1 * 1000 * 1000)
	produceSpeed := time.Duration(speed) * time.Microsecond

	tick := time.Tick(produceSpeed)

	for {
		select {
		case <-tick:
			if vacantSize := atomic.LoadInt64(&b.vacantSize); vacantSize < b.maxSize {
				fmt.Println(vacantSize, b.maxSize)
				atomic.AddInt64(&b.vacantSize, 1)
				go func() {
					b.RWChan <- 1
				}()
			}
		}

	}
}
